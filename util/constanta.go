package util

type ExpenseStatus int32

type ApprovalStatus int32

type UserRole int32

const (
	EXPENSE_PENDING       ExpenseStatus = 3
	EXPENSE_APPROVED      ExpenseStatus = 1
	EXPENSE_REJECTED      ExpenseStatus = -1
	EXPENSE_AUTO_APPROVED ExpenseStatus = 2

	APPROVAL_APPROVED ApprovalStatus = 1
	APPROVAL_REJECTED ApprovalStatus = -1

	USER_ROLE_ADMIN    UserRole = 1
	USER_ROLE_MANAGER  UserRole = 2
	USER_ROLE_EMPLOYEE UserRole = 3
)

func GetExpenseStatusString(status ExpenseStatus) string {
	switch status {
	case EXPENSE_AUTO_APPROVED:
		return "auto_approved"
	case EXPENSE_PENDING:
		return "pending"
	case EXPENSE_APPROVED:
		return "approved"
	case EXPENSE_REJECTED:
		return "rejected"
	}
	return "Unknown"
}

func GetApprovalStatusString(status ApprovalStatus) string {
	switch status {
	case APPROVAL_APPROVED:
		return "approved"
	case APPROVAL_REJECTED:
		return "rejected"
	}
	return "Unknown"
}

func GetUserRoleString(role UserRole) string {
	switch role {
	case USER_ROLE_ADMIN:
		return "admin"
	case USER_ROLE_MANAGER:
		return "manager"
	case USER_ROLE_EMPLOYEE:
		return "employee"
	}
	return "Unknown"
}
