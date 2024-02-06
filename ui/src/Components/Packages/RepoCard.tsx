 

import { ExpandLessOutlined, ExpandMoreOutlined } from '@mui/icons-material'
import { Card, CardActions, CardContent, CardHeader, Collapse, IconButton, Typography } from '@mui/material'
import React, { useState } from 'react'
import { RepositoryType } from '../../api/repository'
import { useAPI } from '../../api/useAPI'
import { curl, aarticlient } from '../../cli/cli'
import { packageTypeIcon } from '../../icons/packageTypeIcon'
import { MultiLangCode, MultiLangCodeItem } from '../MultiLangCode'

export interface RepoCardProps {
  type: RepositoryType
  repo?: string;
  sub?: string;
}

export const RepoCard = ({type, repo = '', sub}: RepoCardProps) => {
  const {credentials} = useAPI();
  const [expanded, setExpanded] = useState(false)
  return (
    <Card>
      <CardHeader avatar={packageTypeIcon(type)} title={repo ? (repo + (sub ? '/' + sub : '')) : sub ? sub : type}
                  titleTypographyProps={{ variant: 'h5' }} />
      <CardContent sx={{pt: 0, pb: 0}}>
        <Typography variant='h6'>Setup</Typography>
        <MultiLangCode storageKey='lang' title='Run this command to setup the repository on your machine :'>
          <MultiLangCodeItem
            label='aarticlient'
            code={aarticlient.setup(type, repo, sub)}
            hiddenCode={aarticlient.setup(type, repo, sub, credentials)}
            language='bash'
          />
          <MultiLangCodeItem
            label='curl'
            code={curl.setup(type, repo, sub)}
            hiddenCode={curl.setup(type, repo, sub, credentials)}
            language='bash'
          />
        </MultiLangCode>
      </CardContent>
      <CardActions sx={{justifyContent: 'end'}}>
        <IconButton onClick={() => setExpanded(!expanded)}>
          {expanded ? <ExpandLessOutlined /> : <ExpandMoreOutlined /> }
        </IconButton>
      </CardActions>
      <Collapse in={expanded} timeout="auto" unmountOnExit>
        <CardContent sx={{pt: 0}}>
          <Typography variant='h6'>Push</Typography>
          <MultiLangCode storageKey='lang' title='Run this command on your machine to push a package to the repository :'>
            <MultiLangCodeItem
              label='aarticlient'
              code={aarticlient.push(type, repo, sub)}
              hiddenCode={aarticlient.push(type, repo, sub, credentials)}
              language='bash'
            />
            <MultiLangCodeItem
              label='curl'
              code={curl.push(type, repo, sub)}
              hiddenCode={curl.push(type, repo, sub, credentials)}
              language='bash'
            />
          </MultiLangCode>
        </CardContent>
      </Collapse>
    </Card>
  )
}
