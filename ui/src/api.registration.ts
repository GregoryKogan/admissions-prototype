import { instance } from '@/api.config'

const RegistrationService = {
  accept: async (registrationId: number) =>
    await instance.post(`/regdata/admin/accept/${registrationId}`),
  reject: async (registrationId: number, reason: string) =>
    await instance.post(`/regdata/admin/reject/${registrationId}`, { reason }),
  list: async () => await instance.get('/regdata/admin/pending'),
  mine: async () => await instance.get('/regdata/mine'),
  register: async (registration: RegistrationRequest) =>
    await instance.post('/regdata', registration),
}

export default RegistrationService

export interface Registration {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  birth_date: string
  email: string
  first_name: string
  gender: string
  grade: number
  june_exam: boolean
  last_name: string
  old_school: string
  parent_first_name: string
  parent_last_name: string
  parent_patronymic: string
  parent_phone: string
  patronymic: string
  source: string
  vmsh: boolean
}

export interface RegistrationRequest {
  birth_date: Date
  email: string
  first_name: string
  gender: string
  grade: number
  june_exam: boolean
  last_name: string
  old_school: string
  parent_first_name: string
  parent_last_name: string
  parent_patronymic: string
  parent_phone: string
  patronymic: string
  source: string
  vmsh: boolean
}
