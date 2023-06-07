import React, { useEffect, useMemo, useState } from 'react'
import {
  Container,
  PageBody,
  Text,
  Color,
  TableV2,
  Layout,
  StringSubstitute,
  Icon,
  FontVariation,
  FlexExpander,
  Utils
} from '@harness/uicore'
import { useHistory } from 'react-router-dom'
import { useGet } from 'restful-react'
import type { CellProps, Column } from 'react-table'
import { Case, Match, Render, Truthy } from 'react-jsx-match'
import ReactTimeago from 'react-timeago'
import { makeDiffRefs, PullRequestFilterOption } from 'utils/GitUtils'
import { useAppContext } from 'AppContext'
import { useGetRepositoryMetadata } from 'hooks/useGetRepositoryMetadata'
import { useStrings } from 'framework/strings'
import { RepositoryPageHeader } from 'components/RepositoryPageHeader/RepositoryPageHeader'
import { voidFn, getErrorMessage, LIST_FETCHING_LIMIT, permissionProps, PageBrowserProps } from 'utils/Utils'
import { usePageIndex } from 'hooks/usePageIndex'
import { useGetSpaceParam } from 'hooks/useGetSpaceParam'
import { useUpdateQueryParams } from 'hooks/useUpdateQueryParams'
import { useQueryParams } from 'hooks/useQueryParams'
import type { TypesPullReq, TypesRepository } from 'services/code'
import { ResourceListingPagination } from 'components/ResourceListingPagination/ResourceListingPagination'
import { UserPreference, useUserPreference } from 'hooks/useUserPreference'
import { NoResultCard } from 'components/NoResultCard/NoResultCard'
import { PipeSeparator } from 'components/PipeSeparator/PipeSeparator'
import { GitRefLink } from 'components/GitRefLink/GitRefLink'
import { PullRequestStateLabel } from 'components/PullRequestStateLabel/PullRequestStateLabel'
import { LoadingSpinner } from 'components/LoadingSpinner/LoadingSpinner'
import { ExecutionStatusLabel } from 'components/ExecutionStatusLabel/ExecutionStatusLabel'
import { PullRequestsContentHeader } from './PullRequestsContentHeader/PullRequestsContentHeader'
import css from './PullRequests.module.scss'

export default function PullRequests() {
  const { getString } = useStrings()
  const history = useHistory()
  const { routes } = useAppContext()
  const [searchTerm, setSearchTerm] = useState('')
  const [filter, setFilter] = useUserPreference<string>(
    UserPreference.PULL_REQUESTS_FILTER_SELECTED_OPTIONS,
    PullRequestFilterOption.OPEN
  )
  const space = useGetSpaceParam()
  const { updateQueryParams } = useUpdateQueryParams()

  const pageBrowser = useQueryParams<PageBrowserProps>()
  const pageInit = pageBrowser.page ? parseInt(pageBrowser.page) : 1
  const [page, setPage] = usePageIndex(pageInit)

  useEffect(() => {
    updateQueryParams({ page: page.toString() })
  }, [setPage]) // eslint-disable-line react-hooks/exhaustive-deps

  const { repoMetadata, error, loading, refetch } = useGetRepositoryMetadata()
  const {
    data,
    error: prError,
    loading: prLoading,
    response
  } = useGet<TypesPullReq[]>({
    path: `/api/v1/repos/${repoMetadata?.path}/+/pullreq`,
    queryParams: {
      limit: String(LIST_FETCHING_LIMIT),
      page,
      sort: filter == PullRequestFilterOption.MERGED ? 'merged' : 'number',
      order: 'desc',
      query: searchTerm,
      state: filter == PullRequestFilterOption.ALL ? '' : filter
    },
    lazy: !repoMetadata
  })

  const { standalone } = useAppContext()
  const { hooks } = useAppContext()
  const permPushResult = hooks?.usePermissionTranslate?.(
    {
      resource: {
        resourceType: 'CODE_REPOSITORY'
      },
      permissions: ['code_repo_push']
    },
    [space]
  )

  const columns: Column<TypesPullReq>[] = useMemo(
    () => [
      {
        id: 'title',
        width: '100%',
        Cell: ({ row }: CellProps<TypesPullReq>) => {
          return (
            <Layout.Horizontal className={css.titleRow} spacing="medium">
              <PullRequestStateLabel iconSize={22} data={row.original} iconOnly />
              <Container padding={{ left: 'small' }}>
                <Layout.Vertical spacing="xsmall">
                  <Text color={Color.GREY_800} className={css.title}>
                    {row.original.title}
                    <Icon
                      className={css.convoIcon}
                      padding={{ left: 'medium', right: 'xsmall' }}
                      name="code-chat"
                      size={15}
                    />
                    <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500} tag="span">
                      {row.original.stats?.conversations}
                    </Text>
                  </Text>
                  <Container>
                    <Layout.Horizontal spacing="small" style={{ alignItems: 'center' }}>
                      <Text color={Color.GREY_500} font={{ size: 'small' }}>
                        <StringSubstitute
                          str={getString('pr.statusLine')}
                          vars={{
                            state: <strong className={css.state}>{row.original.state}</strong>,
                            number: <Text inline>{row.original.number}</Text>,
                            time: (
                              <strong>
                                <ReactTimeago
                                  date={
                                    (row.original.state == 'merged'
                                      ? row.original.merged
                                      : row.original.created) as number
                                  }
                                />
                              </strong>
                            ),
                            user: <strong>{row.original.author?.display_name}</strong>
                          }}
                        />
                      </Text>
                      <PipeSeparator height={10} />
                      <Container>
                        <Layout.Horizontal spacing="xsmall" style={{ alignItems: 'center' }} onClick={Utils.stopEvent}>
                          <GitRefLink
                            text={row.original.target_branch as string}
                            url={routes.toCODERepository({
                              repoPath: repoMetadata?.path as string,
                              gitRef: row.original.target_branch
                            })}
                          />
                          <Text color={Color.GREY_500}>←</Text>
                          <GitRefLink
                            text={row.original.source_branch as string}
                            url={routes.toCODERepository({
                              repoPath: repoMetadata?.path as string,
                              gitRef: row.original.source_branch
                            })}
                          />
                        </Layout.Horizontal>
                      </Container>
                    </Layout.Horizontal>
                  </Container>
                </Layout.Vertical>
              </Container>
              <FlexExpander />
              {/* TODO: Pass proper state when check api is fully implemented */}
              <ExecutionStatusLabel data={{ state: 'success' }} />
            </Layout.Horizontal>
          )
        }
      }
    ],
    [getString] // eslint-disable-line react-hooks/exhaustive-deps
  )

  return (
    <Container className={css.main}>
      <RepositoryPageHeader
        repoMetadata={repoMetadata}
        title={getString('pullRequests')}
        dataTooltipId="repositoryPullRequests"
      />
      <PageBody error={getErrorMessage(error || prError)} retryOnError={voidFn(refetch)}>
        <LoadingSpinner visible={loading} />

        <Render when={repoMetadata}>
          <Layout.Vertical>
            <PullRequestsContentHeader
              loading={prLoading}
              repoMetadata={repoMetadata as TypesRepository}
              activePullRequestFilterOption={filter}
              onPullRequestFilterChanged={_filter => {
                setFilter(_filter)
                setPage(1)
              }}
              onSearchTermChanged={value => {
                setSearchTerm(value)
                setPage(1)
              }}
            />
            <Container padding="xlarge">
              <Match expr={data?.length}>
                <Truthy>
                  <>
                    <TableV2<TypesPullReq>
                      className={css.table}
                      hideHeaders
                      columns={columns}
                      data={data || []}
                      getRowClassName={() => css.row}
                      onRowClick={row => {
                        history.push(
                          routes.toCODEPullRequest({
                            repoPath: repoMetadata?.path as string,
                            pullRequestId: String(row.number)
                          })
                        )
                      }}
                    />
                    <ResourceListingPagination response={response} page={page} setPage={setPage} />
                  </>
                </Truthy>
                <Case val={0}>
                  <NoResultCard
                    forSearch={!!searchTerm}
                    message={getString('pullRequestEmpty')}
                    buttonText={getString('newPullRequest')}
                    onButtonClick={() =>
                      history.push(
                        routes.toCODECompare({
                          repoPath: repoMetadata?.path as string,
                          diffRefs: makeDiffRefs(repoMetadata?.default_branch as string, '')
                        })
                      )
                    }
                    permissionProp={permissionProps(permPushResult, standalone)}
                  />
                </Case>
              </Match>
            </Container>
          </Layout.Vertical>
        </Render>
      </PageBody>
    </Container>
  )
}
