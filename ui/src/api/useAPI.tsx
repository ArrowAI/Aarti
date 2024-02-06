 

import React, { PropsWithChildren, useContext, useEffect, useState } from 'react'
import { useAsync, usePersistedState } from '../hooks'
import { AsyncVoid, Void } from '../utils'
import { API as _API, api } from './api'
import { Credentials } from './schemas/login'

type API = Omit<_API, 'login' | 'logout' | 'credentials'>

interface APIContext extends API {
  login: (credentials: Credentials, repo: string) => Promise<[boolean, Error?]>;
  logout: () => Promise<void>;
  credentials?: Credentials;
  authenticated?: boolean;
  baseRepo?: string;
  setBaseRepo: (repo?: string) => void;
}


const apiContext = React.createContext<APIContext>({
  ...api,
  login: async () => [false],
  logout: AsyncVoid,
  credentials: undefined,
  setBaseRepo: Void,
})


export interface APIProviderProps extends PropsWithChildren<any> {
  user?: string;
  password?: string;
}

export const APIProvider = ({ children }: APIProviderProps) => {
  const [baseRepo, setBaseRepo] = usePersistedState<string | undefined>(undefined, 'baseRepo')
  const [authenticated, setAuthenticated, loaded] = usePersistedState<boolean|undefined>(false, 'authenticated')
  const [credentials, setCredentials] = useState<Credentials>()

  useAsync(async () => {
    if (!authenticated) {
      setCredentials(undefined)
      return
    }
    const [{user, password}, error] = await api.credentials()
    if (error) {
      console.error(error)
      setCredentials(undefined)
      return
    }
    setCredentials({user: user!!, password: password!!})
  }, [authenticated])

  useAsync(async () => {
    if (!loaded || authenticated) return
    // check if auth is required
    const [_, error] = await api.repositories(baseRepo)
    if (error) {
      setAuthenticated(false)
      return
    }
    setAuthenticated(true)
  }, [loaded, authenticated])

  const login = async ({ user, password }: Credentials, repo: string = '') => {
    const [success, error] = await api.login(user, password, repo)
    if (success) setAuthenticated(true)
    return [success, error] as [boolean, Error?]
  }
  const logout = async () => {
    await api.logout()
    setAuthenticated(undefined)
    setBaseRepo(undefined)
  }
  return <apiContext.Provider
    value={{
      ...api,
      credentials,
      login,
      logout,
      authenticated: loaded ? authenticated : undefined,
      baseRepo,
      setBaseRepo,
    }}>{children}</apiContext.Provider>
}

export const useAPI = () => {
  const api = useContext(apiContext)
  return api
}

