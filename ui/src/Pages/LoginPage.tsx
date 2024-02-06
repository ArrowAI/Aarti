 

import { Button, Card, Stack, TextField, Typography } from '@mui/material'
import Box from '@mui/material/Box'
import React, { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { Credentials } from '../api/schemas/login'
import { useAPI } from '../api/useAPI'
import { Loading } from '../Components/Loading'
import { LoginForm } from '../Components/Login/LoginForm'
import aarticlientLogo from '../img/aarticlient-logo.png'
import { MainRoutesRegistry } from '../routes'
import { useSnackbar } from '../snackbar'


export const LoginPage = () => {
  const { errorSnackbar } = useSnackbar()

  const navigate = useNavigate()
  const location = useLocation()
  const from = (location.state as any)?.from?.pathname || '/'

  const [creds, setCreds] = useState<Credentials>()

  const { login: _login, authenticated, baseRepo: repo, setBaseRepo: setRepo, repositories } = useAPI()
  useEffect(() => {
    console.log('authenticated:', authenticated, 'from:', from)
    if (authenticated) navigate(from !== '/login' && from !== '/logout' ? from : '/', { replace: true })
  }, [authenticated])

  const loading = authenticated === undefined

  const [isSubmitting, setIsSubmitting] = useState(false)
  const login = async (c: Credentials) => {
    setIsSubmitting(true)
    const [success] = await _login(c, repo ?? '')
    setIsSubmitting(false)
    if (success) {
      return
    }
    if (repo) {
      setCreds(undefined)
      setRepo(undefined)
      errorSnackbar('invalid username or password')
      return
    }
    setCreds(c)
  }
  const handleSubmit = async (e: any) => {
    setIsSubmitting(true)
    e.preventDefault()
    await login(creds!!)
    setIsSubmitting(false)
  }
  return loading ? <Loading /> : (
    <Stack flex={1} justifyContent='space-around'>
      <Card
        sx={{
          padding: 2,
          paddingTop: 4,
          alignSelf: 'center',
          minWidth: [360, 512],
        }}
      >
        {
          creds
            ? <Stack component='form' onSubmit={handleSubmit}>
              <Stack direction='row'>
                <Box component='img' src={aarticlientLogo}
                     sx={{ height: 20, transform: 'scale(2)', marginLeft: 2, marginRight: 2 }} />
                <Typography variant='h5'>Repository</Typography>
              </Stack>
              <TextField fullWidth={true} autoFocus={true} disabled={isSubmitting}
                         onChange={e => setRepo(e.target.value)} />
              <Button type='submit' disabled={isSubmitting}>Enter</Button>
            </Stack>
            : <>
              <Stack direction='row' paddingBottom={4}>
                <Box component='img' src={aarticlientLogo}
                     sx={{ height: 20, transform: 'scale(2)', marginLeft: 2, marginRight: 2 }} />
                <Typography variant='h5'>Login</Typography>
              </Stack>
              <LoginForm onLogin={login} />
            </>
        }
      </Card>
    </Stack>
  )
}

MainRoutesRegistry['login'] = {
  path: '/login',
  component: <LoginPage />,
  public: true,
  show: false,
  navigate: () => '/login',
}
