import { instance } from '@/api.config'

const RegistrationService = {
  accept: async (registrationId: number) =>
    await instance.post(`/regdata/admin/accept/${registrationId}`),
  reject: async (registrationId: number, reason: string) =>
    await instance.post(`/regdata/admin/reject/${registrationId}`, { reason }),
  list: async () => await instance.get('/regdata/admin/pending'),
}

export default RegistrationService
