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

    const treeSupport = rawData?.repository?.commit?.tree?.codeIntelInfo
    const blobSupport = rawData?.repository?.commit?.blob?.codeIntelSupport

    const data = treeSupport
        ? {
              searchBasedSupport:
                  treeSupport?.searchBasedSupport?.map(support => ({
                      supportLevel: support.support.supportLevel,
                      language: support.support.language || undefined,
                  })) || [],
              preciseSupport:
                  treeSupport?.preciseSupport?.map(support => ({
                      supportLevel: support.support.supportLevel,
                      indexers: support.support.indexers?.map(index => ({ name: index.name, url: index.url })) || [],
                      confidence: support.confidence,
                  })) || [],
          }
        : blobSupport
        ? {
              searchBasedSupport: [
                  {
                      supportLevel: blobSupport.searchBasedSupport.supportLevel,
                      language: blobSupport.searchBasedSupport.language || undefined,
                  },
              ],
              preciseSupport: [
                  {
                      supportLevel: blobSupport.preciseSupport.supportLevel,
                      indexers:
                          blobSupport.preciseSupport.indexers?.map(index => ({
                              name: index.name,
                              url: index.url,
                          })) || undefined,
                  },
              ],
          }
        : undefined

    return {
        data,
        error,
        loading,
    }
}
