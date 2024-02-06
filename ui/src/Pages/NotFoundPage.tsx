 

import { Stack } from '@mui/material'
import { TypingMessage } from '../Components/TypingMessage'
import { NotFoundIcon } from '../icons/NotFoundIcon'
import { MainRoutesRegistry } from '../routes'

export const NotFoundPage = () => (
  <Stack alignItems='center' justifyContent='center' paddingTop={8}>
    <Stack>
      <NotFoundIcon
        sx={{ fontSize: '24rem', padding: theme => theme.spacing(8), color: theme => theme.palette.text.secondary }} />
    </Stack>
    <TypingMessage
      text={[
        '404, Page Not Found.',
        'Sorry... 🙀',
      ]}
      speed={100}
      eraseDelay={3000}
      eraseSpeed={50}
      typingDelay={1000}
    />
  </Stack>
)

MainRoutesRegistry['notFound'] = {
  path: '/*',
  component: <NotFoundPage />,
  priority: 0,
  public: false,
  show: false,
  navigate: () => '/404',
}

