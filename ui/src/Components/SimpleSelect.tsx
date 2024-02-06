 

import { Input, OutlinedInput, Select, SelectProps, TextField, Theme } from '@mui/material'
import { styled } from '@mui/material/styles'

const StyledSelect = styled(Select)(({ theme }) => ({
  border: 'none',
  padding: 0,
  color: theme.palette.primary.main,
  top: 'unset',
  '&.Mui-focused': {
    border: 'none',
    outline: 'none',
  },
  '& > svg': {
    color: theme.palette.primary.main,
  },
}))

export const SimpleSelect = (props: SelectProps) => <StyledSelect {...props} input={<Input />} />
