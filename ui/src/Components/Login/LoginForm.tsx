 

import { LockOutlined, PersonOutlineOutlined, VisibilityOffOutlined, VisibilityOutlined } from '@mui/icons-material'
import { Button, IconButton, InputAdornment, Stack } from '@mui/material'
import { Formik } from 'formik'
import { useState } from 'react'
import { Credentials, credentialsSchema } from '../../api/schemas/login'
import { defaultSpacing } from '../../theme/theme'
import { FTextField } from '../Form'

export const LoginForm = ({ onLogin }: { onLogin: (credentials: Credentials) => Promise<void> }) => {
  const [showPassword, setShowPassword] = useState(false)
  return (
    <Formik<Credentials> initialValues={{ user: '', password: '' }} validationSchema={credentialsSchema}
                         onSubmit={onLogin}>
      {({ isSubmitting, handleReset, handleSubmit }) => (
        <Stack component='form' onReset={handleReset} onSubmit={handleSubmit}>
          <Stack>
            <FTextField
              name='user'
              label='Username'
              autoFocus={true}
              InputProps={{
                startAdornment: (
                  <InputAdornment position='start'>
                    <PersonOutlineOutlined />
                  </InputAdornment>
                ),
              }}
            />
            <FTextField
              name='password'
              label='Password'
              type={showPassword ? 'text' : 'password'}
              InputProps={{
                startAdornment: (
                  <InputAdornment position='start'>
                    <LockOutlined />
                  </InputAdornment>
                ),
                endAdornment: (
                  <InputAdornment position='end'>
                    <IconButton aria-label='toggle password visibility' onClick={() => setShowPassword(!showPassword)}>
                      {showPassword ? <VisibilityOffOutlined /> : <VisibilityOutlined />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          </Stack>
          <Stack direction='row' spacing={defaultSpacing} justifyContent='flex-end'>
            <Button disabled={isSubmitting} type='submit'>
              Login
            </Button>
          </Stack>
        </Stack>
      )}
    </Formik>
  )
}
