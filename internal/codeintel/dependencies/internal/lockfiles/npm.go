package lockfiles

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/sourcegraph/sourcegraph/internal/conf/reposource"
	"github.com/sourcegraph/sourcegraph/internal/lazyregexp"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

//
// package-lock.json
//

type packageLockDependency struct {
	Version      string
	Dev          bool
	Dependencies map[string]*packageLockDependency
}

func parsePackageLockFile(r io.Reader) ([]reposource.PackageVersion, error) {
	var lockFile struct {
		Dependencies map[string]*packageLockDependency
	}

	err := json.NewDecoder(r).Decode(&lockFile)
	if err != nil {
		return nil, errors.Errorf("decode error: %w", err)
	}

	return parsePackageLockDependencies(lockFile.Dependencies)
}

func parsePackageLockDependencies(in map[string]*packageLockDependency) ([]reposource.PackageVersion, error) {
	var (
		errs errors.MultiError
		out  = make([]reposource.PackageVersion, 0, len(in))
	)

	for name, d := range in {
		dep, err := reposource.ParseNpmPackageVersion(name + "@" + d.Version)
		if err != nil {
			errs = errors.Append(errs, err)
		} else {
			out = append(out, dep)
		}

		if d.Dependencies != nil {
			// Recursion
			deps, err := parsePackageLockDependencies(d.Dependencies)
			out = append(out, deps...)
			errs = errors.Append(errs, err)
		}
	}

	return out, errs
}

//
// yarn.lock
//

func parseYarnLockFile(r io.Reader) (deps []reposource.PackageVersion, graph *DependencyGraph, err error) {
	var (
		name string
		skip bool
		errs errors.MultiError

		current             *reposource.NpmPackageVersion
		parsingDependencies bool
	)

	/* yarn.lock

	__metadata:
	  version: 4
	  cacheKey: 6

	"asap@npm:~2.0.6":
	  version: 2.0.6
	  resolution: "asap@npm:2.0.6"
	  checksum: 3d314f8c598b625a98347bacdba609d4c889c616ca5d8ea65acaae8050ab8b7aa6630df2cfe9856c20b260b432adf2ee7a65a1021f268ef70408c70f809e3a39
	  languageName: node
	  linkType: hard
	*/

	var (
		byName          = map[string]*reposource.NpmPackageVersion{}
		dependencyNames = map[*reposource.NpmPackageVersion][]*reposource.NpmPackageVersionConstraint{}
	)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			continue
		}

		var version string
		if version, err = getVersion(line); err == nil { // e.g. version: 2.0.6
			if skip {
				continue
			}

			if name == "" {
				return nil, nil, errors.New("invalid yarn.lock format")
			}

			dep, err := reposource.ParseNpmPackageVersion(name + "@" + version)
			if err != nil {
				errs = errors.Append(errs, err)
			} else {
				deps = append(deps, dep)
				byName[name] = dep
				current = dep
				dependencyNames[current] = []*reposource.NpmPackageVersionConstraint{}
				name = ""
			}
			continue
		}

		if skip = strings.HasPrefix(line, "__metadata"); skip {
			continue
		}

		if line[:1] != " " && line[:1] != "#" { // e.g. "asap@npm:~2.0.6":
			parsingDependencies = false

			var packagename string
			var protocols []string
			if packagename, protocols, err = parsePackageLocatorLine(line); err != nil {
				continue
			}
			fmt.Printf("packagename=%q, protocols=%+v\n", packagename, protocols)
			skipOuter := false
			for _, protocol := range protocols {
				if skip = !validProtocol(protocol); skip {
					skipOuter = true
				}
			}
			if skipOuter {
				continue
			}
			name = packagename
			current = nil
		}

		if line == "  dependencies:" {
			parsingDependencies = true
		}

		if line[:4] == "    " && parsingDependencies && current != nil {
			depName, versionConstraint, err := parsePackageDependency(line)
			if err != nil {
				continue
			}

			fmt.Printf("depName=%s, depVersion=%s\n", depName, versionConstraint)
			dep, err := reposource.ParseNpmPackageVersionConstraint(depName, versionConstraint)
			if err != nil {
				errs = errors.Append(errs, err)
			}

			if deps, ok := dependencyNames[current]; !ok {
				dependencyNames[current] = []*reposource.NpmPackageVersionConstraint{dep}
			} else {
				dependencyNames[current] = append(deps, dep)
			}
		}
	}

	graph = newDependencyGraph()
	for pkg, deps := range dependencyNames {
		graph.addPackage(pkg)

		for _, dep := range deps {
			fmt.Printf("dep = %+v\n", dep)
			// graph.addDependency(pkg, dep)
		}
	}

	return deps, graph, errs
}

var (
	yarnLocatorRegexp    = lazyregexp.New(`"?(?P<package>.+?)@(?P<protocol>[^"]+)"?`)
	yarnDependencyRegexp = lazyregexp.New(`\s{4}"?(?P<package>.+?)"?\s"?(?P<version>[^"]+)"?`)
	yarnVersionRegexp    = lazyregexp.New(`\s+"?version:?"?\s+"?(?P<version>[^"]+)"?`)
)

func parsePackageLocatorLine(target string) (packagename string, protocols []string, err error) {
	trimmed := strings.TrimSuffix(target, ":")
	elems := strings.Split(trimmed, ", ")

	for _, elem := range elems {
		capture := yarnLocatorRegexp.FindStringSubmatch(elem)
		if len(capture) < 2 {
			return "", protocols, errors.New("not package format")
		}
		for i, group := range yarnLocatorRegexp.SubexpNames() {
			switch group {
			case "package":
				packagename = capture[i]
			case "protocol":
				protocols = append(protocols, capture[i])
			}
		}
	}
	return
}

func parsePackageDependency(target string) (dependencyname, version string, err error) {
	capture := yarnDependencyRegexp.FindStringSubmatch(target)
	if len(capture) < 2 {
		return "", "", errors.New("not package format")
	}
	for i, group := range yarnDependencyRegexp.SubexpNames() {
		switch group {
		case "package":
			dependencyname = capture[i]
		case "version":
			version = capture[i]
		}
	}
	return
}

func getVersion(target string) (version string, err error) {
	capture := yarnVersionRegexp.FindStringSubmatch(target)
	if len(capture) < 2 {
		return "", errors.New("not version")
	}
	return capture[len(capture)-1], nil
}

func validProtocol(protocol string) (valid bool) {
	switch protocol {
	// only scan npm packages
	case "npm", "":
		return true
	}
	return false
}
