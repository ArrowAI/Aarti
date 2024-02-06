 

import { Check, ContentCopyOutlined } from '@mui/icons-material'
import { Box, IconButton, Stack, SxProps } from '@mui/material'
import copy from 'copy-to-clipboard'
import Highlight, { defaultProps, Language } from 'prism-react-renderer'
import { useState } from 'react'
import '../xonokai.css'

export interface CodeProps {
  code: string;
  // hiddenCode is used to display a different code than the one copied to the clipboard
  hiddenCode?: string;
  language?: Language;
  sx?: SxProps;
}

export const Code = ({ code, hiddenCode, language, sx }: CodeProps) => {
  const [copied, setCopied] = useState(false)
  const copyToClipboard = async () => {
    copy(hiddenCode || code, { format: 'text/plain', onCopy: () => setCopied(true) })
    setCopied(true)
  }
  return (
    <Box>
      <Stack component='code' direction='row' sx={{ backgroundColor: 'black' }}>
        {
          language ? (
              <Highlight {...defaultProps} theme={undefined} code={code} language={language}>
                {({ className, style, tokens, getLineProps, getTokenProps }) => (
                  <Box
                    className={className}
                    style={style}
                    component='pre'
                    flex={1}
                    sx={{
                      ...sx,
                      margin: 0,
                      padding: 2,
                      backgroundColor: 'black',
                      color: 'primary.main',
                      overflow: 'scroll',
                    }}
                  >
                    {tokens.map((line, i) => (
                      <Box {...getLineProps({ line, key: i })}>
                        {line.map((token, key) => (
                          <Box component='span' {...getTokenProps({ token, key })} />
                        ))}
                      </Box>
                    ))}
                  </Box>
                )}
              </Highlight>
            )
            : (
              <Box
                component='pre'
                flex={1}
                sx={{ padding: 2, margin: 0, backgroundColor: 'black', color: 'primary.main', overflow: 'scroll' }}
              >
                {code}
              </Box>
            )
        }
        <IconButton
          sx={{
            color: 'white',
            opacity: 0.5,
            '&:hover': { opacity: 1 },
          }}
          onMouseLeave={() => setCopied(false)}
          onClick={copyToClipboard}
        >
          {(copied ? <Check /> : <ContentCopyOutlined />)}
        </IconButton>
      </Stack>
    </Box>
  )
}
