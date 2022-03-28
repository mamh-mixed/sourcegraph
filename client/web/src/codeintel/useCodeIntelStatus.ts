import { ApolloError } from '@apollo/client'

import { gql, useQuery } from '@sourcegraph/http-client'

import {
    TreeAndBlobCodeIntelStatusVariables,
    TreeAndBlobCodeIntelStatusResult,
    PreciseSupportLevel,
    InferedPreciseSupportLevel,
    SearchBasedSupportLevel,
} from '../graphql-operations'

const BLOB_AND_TREE_CODE_INTEL_STATUS_QUERY = gql`
    query TreeAndBlobCodeIntelStatus($repository: String!, $commit: String!, $path: String!) {
        repository(name: $repository) {
            commit(rev: $commit) {
                tree(path: $path) {
                    codeIntelInfo {
                        searchBasedSupport {
                            support {
                                supportLevel
                                language
                            }
                        }
                        preciseSupport {
                            support {
                                supportLevel
                                indexers {
                                    name
                                    url
                                }
                            }
                            confidence
                        }
                    }
                    lsif {
                        lsifUploads {
                            id
                        }
                    }
                }
                blob(path: $path) {
                    codeIntelSupport {
                        searchBasedSupport {
                            supportLevel
                            language
                        }
                        preciseSupport {
                            supportLevel
                            indexers {
                                name
                                url
                            }
                        }
                    }
                    lsif {
                        lsifUploads {
                            id
                        }
                    }
                }
            }
        }
    }
`

interface UseCodeIntelStatusParameters {
    variables: TreeAndBlobCodeIntelStatusVariables
}

interface UseCodeIntelStatusResult {
    data?: {
        searchBasedSupport: {
            supportLevel: SearchBasedSupportLevel
            language?: string
        }[]
        preciseSupport: {
            supportLevel: PreciseSupportLevel
            indexers?: { name: string; url: string }[]
            confidence?: InferedPreciseSupportLevel
        }[]
        uploadIds: string[]
    }
    error?: ApolloError
    loading: boolean
}

export const useCodeIntelStatus = ({ variables }: UseCodeIntelStatusParameters): UseCodeIntelStatusResult => {
    const { data: rawData, error, loading } = useQuery<
        TreeAndBlobCodeIntelStatusResult,
        TreeAndBlobCodeIntelStatusVariables
    >(BLOB_AND_TREE_CODE_INTEL_STATUS_QUERY, {
        variables,
        notifyOnNetworkStatusChange: false,
        fetchPolicy: 'no-cache',
        errorPolicy: 'ignore', // TODO - necessary because tree OR blob will fail
    })

    const tree = rawData?.repository?.commit?.tree
    const blob = rawData?.repository?.commit?.blob

    const data = tree
        ? {
              searchBasedSupport:
                  tree.codeIntelInfo?.searchBasedSupport?.map(support => ({
                      supportLevel: support.support.supportLevel,
                      language: support.support.language || undefined,
                  })) || [],
              preciseSupport:
                  tree.codeIntelInfo?.preciseSupport?.map(support => ({
                      supportLevel: support.support.supportLevel,
                      indexers: support.support.indexers?.map(index => ({ name: index.name, url: index.url })) || [],
                      confidence: support.confidence,
                  })) || [],
              uploadIds: tree.lsif?.lsifUploads.map(upload => upload.id) || [],
          }
        : blob
        ? {
              searchBasedSupport: [
                  {
                      supportLevel: blob.codeIntelSupport.searchBasedSupport.supportLevel,
                      language: blob.codeIntelSupport.searchBasedSupport.language || undefined,
                  },
              ],
              preciseSupport: [
                  {
                      supportLevel: blob.codeIntelSupport.preciseSupport.supportLevel,
                      indexers:
                          blob.codeIntelSupport.preciseSupport.indexers?.map(index => ({
                              name: index.name,
                              url: index.url,
                          })) || undefined,
                  },
              ],
              uploadIds: blob.lsif?.lsifUploads.map(upload => upload.id) || [],
          }
        : undefined

    return {
        data,
        error,
        loading,
    }
}
