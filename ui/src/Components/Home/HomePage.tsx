 

import { Box, Card, CardContent, CardHeader, Stack, Typography } from '@mui/material'
import React, { useState } from 'react'
import { Repository } from '../../api/repository'
import { useAPI } from '../../api/useAPI'
import { aarticlient } from '../../cli/cli'
import { useAsync } from '../../hooks'
import { useSnackbar } from '../../snackbar'

import { defaultPadding, defaultSpacing } from '../../theme/theme'
import { Loading } from '../Loading'
import { MultiLangCode, MultiLangCodeItem } from '../MultiLangCode'
import { RepositoryCard } from './RepositoryCard'

const HomePage = () => {
  const api = useAPI()
  const [repos, setRepos] = useState<Repository[]>([])
  const { errorSnackbar } = useSnackbar()
  const [loading, setLoading] = useState(false)
  useAsync(async () => {
    setLoading(true)
    const [repos, error] = await api.repositories(api.baseRepo)
    setRepos(repos)
    if (error) {
      errorSnackbar(error.message)
    }
    setLoading(false)
  }, [])
  return loading ? <Loading /> : (
    <Stack padding={defaultPadding} spacing={defaultSpacing}>
      <Stack direction='row'>
        <Box flex={1}>
          <Typography variant='h6'>Repositories</Typography>
        </Box>
      </Stack>
      <Stack>
        <Card>
          {
            api.baseRepo
            && <CardHeader
              title={api.baseRepo}
              titleTypographyProps={{ variant: 'h5' }}
            />
          }
          <CardContent>
            <MultiLangCode storageKey='lang' title='Run this command to log into the repository on your machine :'>
              <MultiLangCodeItem
                label='aarticlient'
                code={aarticlient.login(api.baseRepo)}
                hiddenCode={aarticlient.login(api.baseRepo, api.credentials)}
                language='bash'
              />
            </MultiLangCode>
          </CardContent>
        </Card>
      </Stack>
      <Stack>
        {
          repos.map(({ name, type, ...rest }) => (
            <RepositoryCard key={`${name}:${type}`} repository={{ name, type, ...rest }} />
          ))
        }
      </Stack>
    </Stack>
  )
}

export default HomePage
