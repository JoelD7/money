export type Expense = {
    expenseID: string
    username: string
    categoryID?: string
    categoryName?: string
    amount: number
    name: string
    notes?: string
    createdDate: Date
    period: string
    updateDate: Date
}