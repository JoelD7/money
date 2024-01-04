export type Expense = {
    expenseID: string
    username: string
    categoryID?: string
    categoryName?: string
    amount: number
    name: string
    createdDate: Date
    period: string
    updateDate: Date
}