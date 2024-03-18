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

export type Category = {
    categoryID: string
    name: string
    budget: number
    color: string
}

export type SignUpUser = {
    username: string
    password: string
    fullname: string
}

export type LoginCredentials = {
    username: string
    password: string
}