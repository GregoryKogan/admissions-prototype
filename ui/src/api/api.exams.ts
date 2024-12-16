import { instance } from './api.config'

const ExamsService = {
  list: async () => await instance.get('/exams/admin'),
  create: async (exam: ExamRequest) =>
    await instance.post('/exams/admin', exam),
  delete: async (examId: number) =>
    await instance.delete(`/exams/admin/${examId}`),
  types: async () => await instance.get('/exams/admin/types'),
  allocation: async (examId: number) =>
    await instance.get(`/exams/allocation/${examId}`),
  available: async () => await instance.get('/exams/available'),
  register: async (examId: number) =>
    await instance.post(`/exams/register/${examId}`),
  registrationStatus: async (examId: number) =>
    await instance.get(`/exams/registration_status/${examId}`),
  history: async () => await instance.get('/exams/history'),
  downloadRegistrations: async (examId: number) =>
    await instance.get(`/exams/admin/registrations/${examId}/download`, {
      responseType: 'blob',
    }),
}

export default ExamsService

export interface Exam {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  start: string
  end: string
  location: string
  capacity: number
  grade: number
  type: ExamType
}

export interface ExamRequest {
  start: Date
  end: Date
  location: string
  capacity: number
  grade: number
  type_id: number
}

export interface ExamType {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  title: string
  order: number
  dismissing: boolean
  has_points: boolean
}

export interface Allocation {
  capacity: number
  occupied: number
}
