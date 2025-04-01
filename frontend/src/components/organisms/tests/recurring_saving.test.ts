// Mock the Date object
import { afterAll, describe, expect, test } from "@jest/globals";
import { SavingGoal } from "../../../types";
import { estimateSavingAmount, getMonthDifference } from "../../../utils/functions.ts";

// Mock current date for testing
let mockCurrentDate: Date;

// Save original Date constructor
const OriginalDate = global.Date;

// Mock Date constructor
global.Date = class extends OriginalDate {
    constructor(...args: any[]) {
        if (args.length === 0) {
            super(mockCurrentDate);
            return;
        }
        // @ts-expect-error ...
        super(...args);
    }

    // Mock Date.now()
    static now() {
        return mockCurrentDate ? mockCurrentDate.getTime() : OriginalDate.now();
    }
} as DateConstructor;

// Function to set the mock date for testing
function setMockDate(date: Date): void {
    mockCurrentDate = date;
}

describe('estimateSavingAmount', () => {

    // Test case 1: savingGoal is undefined
    test('returns 0 when savingGoal is undefined', () => {
        const result = estimateSavingAmount(undefined);

        expect(result).toBe(0);
    });

    // Test case 2: Same month and year
    test('calculates correctly when deadline is in the same month and year', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1000,
            progress: 0,
            deadline: '2025-03-31', // March 31, 2025
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        expect(result).toBe(1000); // Full amount in one month
    });

    // Test case 3: Different month, same year
    test('calculates correctly when deadline is in a different month, same year', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1000,
            progress: 0,
            deadline: '2025-06-15', // June 15, 2025
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        expect(result).toBe(333.4); // 1000 / 3 months = 333.33, rounded up to nearest 0.1
    });

    // Test case 4: Different year
    test('calculates correctly when deadline is in a different year', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1200,
            progress: 0,
            deadline: '2026-03-15', // March 15, 2026
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        // 12 months, plus msInAMonth adjustment for different year
        expect(result).toBe(100); // 1200 / 12 = 100
    });

    // Test case 5: With existing progress
    test('considers existing progress in calculation', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1000,
            progress: 400, // Already saved 400
            deadline: '2025-06-15', // June 15, 2025
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        expect(result).toBe(200); // (1000 - 400) / 3 months = 200
    });

    // Test case 6: Target exactly reached
    test('handles case where progress equals target', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1000,
            progress: 1000, // Already saved the full amount
            deadline: '2025-06-15', // June 15, 2025
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        expect(result).toBe(0); // (1000 - 1000) / 3 months = 0
    });

    // Test case 7: Edge case - deadline is next month
    test('calculates correctly when deadline is next month', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1000,
            progress: 250,
            deadline: '2025-04-01', // April 1, 2025
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        expect(result).toBe(750); // (1000 - 250) / 1 month = 750
    });

    // Test case 8: Edge case - deadline in same month (monthsUntilDeadline = 0)
    test('handles deadline in same month resulting in zero months', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1000,
            progress: 500,
            deadline: '2025-03-14', // March 14, 2025 (already passed)
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);
        expect(result).toBe(savingGoal.target - savingGoal.progress);
    });

    // Test case 9: Long-term goal (several years)
    test('calculates correctly for long-term goals spanning multiple years', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Retirement Fund',
            target: 60000,
            progress: 0,
            deadline: '2030-03-15', // 5 years later
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        // 60 months (5 years) with adjustment for years
        expect(result).toBe(1000); // 60000 / 60 = 1000
    });

    test('calculates correctly for long-term goals spanning multiple years', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Retirement Fund',
            target: 89543,
            progress: 1289,
            deadline: '2033-09-25',
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        expect(result).toBe(865.3); //(89543 - 1289) / 102 = 865.3
    });

    // Test case 10: Rounding behavior
    test('correctly rounds to the nearest 0.1', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1003,
            progress: 0,
            deadline: '2025-07-15', // July 15, 2025 (4 months)
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        expect(result).toBe(250.8);
    });

    // Test case 11: Deadline in the past
    test('handles deadline in the past (negative months)', () => {
        setMockDate(new Date(2025, 5, 15)); // June 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1000,
            progress: 500,
            deadline: '2025-03-15', // March 15, 2025 (past)
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        // monthsUntilDeadline will be negative, so we should expect a negative result
        expect(result).toBeLessThan(0);
    });

    // Test case 12: Edge case - progress greater than target
    test('handles progress greater than target', () => {
        setMockDate(new Date(2025, 2, 15)); // March 15, 2025

        const savingGoal: SavingGoal = {
            saving_goal_id: '123',
            name: 'Test Goal',
            target: 1000,
            progress: 1200, // Already saved more than the target
            deadline: '2025-06-15', // June 15, 2025
            username: 'testuser'
        };

        const result = estimateSavingAmount(savingGoal);

        expect(result).toBeLessThan(0); // (1000 - 1200) / 3 months = -66.7
    });

    // Restore the original Date after all tests
    afterAll(() => {
        global.Date = OriginalDate;
    });
});

describe('getMonthDifference', () => {
    test('returns 0 when dates are in the same month and year', () => {
        const current = new Date("2025-03-01");
        const deadline = new Date("2025-03-15");
        expect(getMonthDifference(current, deadline)).toBe(0);
    });

    test('calculates correctly when dates are in different months of the same year', () => {
        const current = new Date("2025-03-31");
        const deadline = new Date("2025-06-01");
        expect(getMonthDifference(current, deadline)).toBe(3);
    });

    test('calculates correctly when dates are in different years', () => {
        const current = new Date("2025-05-31");
        const deadline = new Date("2026-03-01");
        expect(getMonthDifference(current, deadline)).toBe(10);
    });

    test('calculates correctly when dates are in different months and years', () => {
        const current = new Date("2025-03-15");
        const deadline = new Date("2026-06-15");
        expect(getMonthDifference(current, deadline)).toBe(15);
    });

    test('calculates correctly when dates are in the same month of different years', () => {
        const current = new Date("2025-06-15");
        const deadline = new Date("2027-06-15");
        expect(getMonthDifference(current, deadline)).toBe(24);
    });

    test('handles deadline month before current month in different year', () => {
        const current = new Date("2025-09-15");
        const deadline = new Date("2026-03-15");
        expect(getMonthDifference(current, deadline)).toBe(6);
    });

    test('calculates correctly with multiple years difference', () => {
        const current = new Date("2025-03-15");
        const deadline = new Date("2030-06-15");
        expect(getMonthDifference(current, deadline)).toBe(63);
    });

    test('returns negative months when deadline is before current date', () => {
        const current = new Date("2025-06-15");
        const deadline = new Date("2025-03-15");
        expect(getMonthDifference(current, deadline)).toBe(-3);
    });

    test('calculates one month difference correctly', () => {
        const current = new Date("2025-03-15");
        const deadline = new Date("2025-04-15");
        expect(getMonthDifference(current, deadline)).toBe(1);
    });

    test('calculates correctly from end of month to beginning of next', () => {
        const current = new Date("2025-03-31");
        const deadline = new Date("2025-04-01");
        expect(getMonthDifference(current, deadline)).toBe(1);
    });

    test('handles deadline year before current year', () => {
        const current = new Date("2025-03-15");
        const deadline = new Date("2024-06-15");
        expect(getMonthDifference(current, deadline)).toBe(-9);
    });

    test('works correctly across leap years', () => {
        const current = new Date("2024-02-29");
        const deadline = new Date("2025-02-28");
        expect(getMonthDifference(current, deadline)).toBe(12);
    });
});
