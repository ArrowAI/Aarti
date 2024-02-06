 


import { Tooltip as MuiTooltip, TooltipProps as MuiTooltipProps } from '@mui/material'

export interface TooltipProps extends MuiTooltipProps {
  disabled?: boolean
}

export const Tooltip = ({ disabled, children, ...props }: TooltipProps) => disabled ? <>{children}</> :
  <MuiTooltip {...props}>{children}</MuiTooltip>
