import { ApolloError, ApolloQueryResult, gql, MutationFunctionOptions, FetchResult, makeVar } from '@apollo/client'

import { useQuery, useMutation } from '@sourcegraph/http-client'

import {
    Maybe,
    Exact,
    UserRepositoriesResult,
    UserRepositoriesVariables,
    SetExternalServiceReposResult,
    SetExternalServiceReposVariables,
} from '../graphql-operations'

export interface MinSelectedRepo {
    name: string
    externalRepository: {
        serviceType: string
        id?: string
    }
}

type SelectedRepos = MinSelectedRepo[] | undefined
const SelectedReposInitialValue: SelectedRepos = undefined

// eslint-disable-next-line unicorn/prevent-abbreviations
export const selectedReposVar = makeVar<SelectedRepos>(SelectedReposInitialValue)

interface UseSelectedReposResult {
    selectedRepos:
        | (NonNullable<UserRepositoriesResult['node']> & { __typename: 'User' })['repositories']['nodes']
        | undefined
    loadingSelectedRepos: boolean
    errorSelectedRepos: ApolloError | undefined
    refetchSelectedRepos:
        | ((
              variables?:
                  | Partial<
                        Exact<{
                            user: string
                            codeHost: Maybe<string>
                            query: Maybe<string>
                        }>
                    >
                  | undefined
          ) => Promise<ApolloQueryResult<UserRepositoriesResult>>)
        | undefined
}

const SAVE_SELECTED_REPOS = gql`
    mutation SetExternalServiceRepos($id: ID!, $allRepos: Boolean!, $repos: [String!]) {
        setExternalServiceRepos(id: $id, allRepos: $allRepos, repos: $repos) {
            alwaysNil
        }
    }
`

type UseSaveSelectedReposResult = (
    options?:
        | MutationFunctionOptions<
              SetExternalServiceReposResult,
              Exact<{
                  id: string
                  allRepos: boolean
                  repos: Maybe<string[]>
              }>
          >
        | undefined
) => Promise<FetchResult<SetExternalServiceReposResult>>

export const useSaveSelectedRepos = (): UseSaveSelectedReposResult => {
    const [saveSelectedRepos] = useMutation<SetExternalServiceReposResult, SetExternalServiceReposVariables>(
        SAVE_SELECTED_REPOS
    )

    return saveSelectedRepos
}

export const SELECTED_REPOS = gql`
    query UserSelectedRepositories(
        $id: ID!
        $first: Int
        $query: String
        $cloned: Boolean
        $notCloned: Boolean
        $indexed: Boolean
        $notIndexed: Boolean
        $externalServiceID: ID
    ) {
        node(id: $id) {
            ... on User {
                __typename
                repositories(
                    first: $first
                    query: $query
                    cloned: $cloned
                    notCloned: $notCloned
                    indexed: $indexed
                    notIndexed: $notIndexed
                    externalServiceID: $externalServiceID
                ) {
                    nodes {
                        id
                        name
                        createdAt
                        viewerCanAdminister
                        url
                        isPrivate
                        mirrorInfo {
                            cloned
                            cloneInProgress
                            cloneProgress
                            updatedAt
                        }
                        externalRepository {
                            serviceType
                            serviceID
                        }
                    }
                    totalCount(precise: true)
                    pageInfo {
                        hasNextPage
                    }
                }
            }
        }
    }
`

export const useSelectedRepos = (userId: string, first?: number): UseSelectedReposResult => {
    const { data, loading, error, refetch } = useQuery<UserRepositoriesResult, UserRepositoriesVariables>(
        SELECTED_REPOS,
        {
            variables: {
                id: userId,
                cloned: true,
                notCloned: true,
                indexed: true,
                notIndexed: true,
                first: first || 2000,
                query: null,
                externalServiceID: null,
            },
        }
    )

    return {
        selectedRepos: (data?.node?.__typename === 'User' && data.node.repositories.nodes) || undefined,
        loadingSelectedRepos: loading,
        errorSelectedRepos: error,
        refetchSelectedRepos: refetch,
    }
}
