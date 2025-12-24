package storage

// implemented as NewUserStore and similar functions in concrete store implementations
// different in sqlite or other db implementations
type UserStoreProvider func(exec SQLExecutor) UserStore
type FamilyStoreProvider func(exec SQLExecutor) FamilyStore
type MembershipStoreProvider func(exec SQLExecutor) MembershipStore
type RefreshTokenStoreProvider func(exec SQLExecutor) RefreshTokenStore
