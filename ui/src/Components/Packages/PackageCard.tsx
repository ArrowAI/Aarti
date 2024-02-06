 

import { ExpandLessOutlined, ExpandMoreOutlined, MemoryOutlined } from '@mui/icons-material'
import {
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Chip,
  Collapse,
  IconButton,
  Stack,
  Typography,
} from '@mui/material'
import React, { useState } from 'react'
import { Package, RepositoryType } from '../../api/repository'
import { useAPI } from '../../api/useAPI'
import { curl, aarticlient } from '../../cli/cli'
import { BalanceIcon } from '../../icons/BalanceIcon'
import { KubernetesIcon } from '../../icons/KubernetesIcon'
import { LinuxIcon } from '../../icons/LinuxIcon'
import { packageTypeIcon } from '../../icons/packageTypeIcon'
import { VersionIcon } from '../../icons/VersionIcon'
import { defaultPadding, defaultSpacing } from '../../theme/theme'
import { humanSize } from '../../utils'
import { ExternalLink } from '../ExternalLink'
import { MultiLangCode, MultiLangCodeItem } from '../MultiLangCode'

export interface PackageCardProps {
  repo?: string
  package: Package
}

export const PackageCard = ({
                              repo,
                              package: {
                                name,
                                type,
                                size,
                                version,
                                architecture,
                                license,
                                projectURL,
                                description,
                                filePath,
                              },
                            }: PackageCardProps) => {
  const { credentials } = useAPI()
  const [expanded, setExpanded] = useState(false)
  return (
    <Card>
      <CardHeader
        avatar={packageTypeIcon(type)}
        title={name} subheader={humanSize(size)}
        action={(
          <Stack direction='row' padding={defaultPadding} alignItems='center'>
            <VersionIcon />
            <Typography
              sx={{ marginLeft: '4px !important' }}
              variant='body2'>{version}</Typography>
          </Stack>
        )} />
      <CardContent sx={{ pt: 0, pb: 0 }}>
        <Stack direction='row' marginTop={0}>
          <Chip icon={type === RepositoryType.HELM ? <KubernetesIcon /> : <LinuxIcon sx={{ padding: 0.15 }} />}
                label={type === RepositoryType.HELM ? 'kubernetes' : 'linux'} />
          <Chip icon={<MemoryOutlined />} label={architecture} />
          {license && <Chip icon={<BalanceIcon sx={{ padding: 0.25 }} />} label={license} />}
        </Stack>
        <Stack sx={{ marginTop: defaultSpacing }}>
          <Typography variant='body2' fontStyle='italic'>{description}</Typography>
        </Stack>
      </CardContent>
      <CardActions sx={{ justifyContent: projectURL ? 'space-between' : 'end' }}>
        {projectURL && <ExternalLink href={projectURL}>{projectURL}</ExternalLink>}
        <IconButton onClick={() => setExpanded(!expanded)}>
          {expanded ? <ExpandLessOutlined /> : <ExpandMoreOutlined />}
        </IconButton>
      </CardActions>
      <Collapse in={expanded} timeout='auto' unmountOnExit>
        <CardContent sx={{ pt: 0 }}>
          <Typography variant='h6'>Delete</Typography>
          <MultiLangCode storageKey='lang'
                         title='Run this command on your machine to delete the package from the repository:'>
            <MultiLangCodeItem
              label='aarticlient'
              code={aarticlient.delete(type, filePath, repo)}
              hiddenCode={aarticlient.delete(type, filePath, repo, credentials)}
              language='bash'
            />
            <MultiLangCodeItem
              label='curl'
              code={curl.delete(type, filePath, repo)}
              hiddenCode={curl.delete(type, filePath, repo, credentials)}
              language='bash'
            />
          </MultiLangCode>
        </CardContent>
      </Collapse>
    </Card>
  )
}
