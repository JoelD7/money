export type User = {
    username: string
    currentPeriod: string
    remainder: number
    expenses?:number
    categories?: Category[]
}

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
    id: string
    name: string
    budget: number
    color: string
}

export type SignUpUser = {
    username: string
    password: string
    fullname: string
}

export type Credentials = {
    username: string
    password: string
}