 

import { object, SchemaOf, string } from 'yup'

export interface Credentials {
  user: string
  password: string
}

export const credentialsSchema: SchemaOf<Credentials> = object({
  user: string().min(1).required('Required'),
  password: string().min(1).required('Required'),
})
