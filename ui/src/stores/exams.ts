import { defineStore } from 'pinia'
import ExamsService, { Exam, Allocation } from '@/api/api.exams'

export const useExamsStore = defineStore('exams', {
  state: () => ({
    availableExams: [] as Exam[],
    historyExams: [] as Exam[],
    allocations: new Map<number, Allocation>(),
    registrationStatuses: new Map<
      number,
      { registered: boolean; registered_to_same_type: boolean }
    >(),
  }),
  actions: {
    async reloadAvailableExams() {
      const response = await ExamsService.available()
      this.availableExams = response.data
    },
    async reloadHistoryExams() {
      const response = await ExamsService.history()
      this.historyExams = response.data
    },
    async reloadAll() {
      await Promise.all([
        this.reloadAvailableExams(),
        this.reloadHistoryExams(),
      ])
    },
    async reloadExamData(examId: number) {
      const [allocation, status] = await Promise.all([
        ExamsService.allocation(examId),
        ExamsService.registrationStatus(examId),
      ])
      this.allocations.set(examId, allocation.data)
      this.registrationStatuses.set(examId, status.data)
    },

    async reloadAllExamData() {
      const allExams = [...this.availableExams, ...this.historyExams]
      await Promise.all(allExams.map((exam) => this.reloadExamData(exam.ID)))
    },
  },
})
