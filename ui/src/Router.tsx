 

import React, { PropsWithChildren, useCallback } from 'react'
import { Navigate, Route, Routes, useLocation, useParams } from 'react-router-dom'
import { useAPI } from './api/useAPI'
import { Loading } from './Components/Loading'
import { MainRoutesRegistry, RouteDefinition, RoutesRegistry } from './routes'

export const Router = () => {
  const makeRoutes = useCallback(
    (reg?: RoutesRegistry, parent?: RouteDefinition) =>
      Object.values(reg ?? {}).map((r, i) => {
        return r.public ? (
          <Route key={i} path={parent ? parent.path + '/' + r.path : r.path} element={r.component}>
            {makeRoutes(r.subRoutes, r)}
          </Route>
        ) : (
          <Route key={i} path={parent ? parent.path + '/' + r.path : r.path}
                 element={<ProtectedRoute>{r.component}</ProtectedRoute>}>
            {makeRoutes(r.subRoutes, r)}
          </Route>
        )
      }),
    [],
  )
  return <Routes>{makeRoutes(MainRoutesRegistry)}</Routes>
}

export const ProtectedRoute = ({ children }: PropsWithChildren<any>) => {
  const from = useLocation()
  const fromParams = useParams()
  const { authenticated } = useAPI()
  if (authenticated === undefined) {
    return <Loading />
  }
  if (!authenticated) {
    return <Navigate to={MainRoutesRegistry.login.navigate()} state={{ from, fromParams }} replace />
  }
  return children
}
