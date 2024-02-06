 


import React from 'react'

export const Edit = 'edit'
export const New = 'new'

export interface RouteDefinition {
  priority?: number;
  path: string;
  label?: string;
  component?: React.ReactNode;
  show?: boolean;
  icon?: React.ReactElement;
  hasBottomNavigation?: boolean,
  bottomEnd?: boolean
  subRoutes?: RoutesRegistry;
  public?: boolean
  navigate: (args?: any) => string
}

export interface RoutesRegistry {
  [key: string]: RouteDefinition;
}

export const MainRoutesRegistry: RoutesRegistry = {}
