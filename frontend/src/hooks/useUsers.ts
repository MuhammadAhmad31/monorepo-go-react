import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useListUsers,
  useCreateUser,
  useUpdateUser,
  useDeleteUser,
  useGetUser,
  getListUsersQueryKey,
  getGetUserQueryKey,
} from '../generated/users/users'
import type {
  CreateUserRequest,
  UpdateUserRequest,
  ListUsersParams,
} from '../generated/models'

const useUsersCache = () => {
  const queryClient = useQueryClient()

  const invalidateUsersList = () =>
    queryClient.invalidateQueries({
      queryKey: getListUsersQueryKey(),
      exact: false,
    })

  return { queryClient, invalidateUsersList }
}

export const useUsersList = (params?: ListUsersParams) => {
  const query = useListUsers(params)

  return {
    users: query.data?.data ?? [],
    meta: query.data?.meta,
    isLoading: query.isLoading,
    isError: query.isError,
    error: query.error,
    refetch: query.refetch,
  }
}

export const useUserDetail = (id: string) => {
  const query = useGetUser(id)

  return {
    user: query.data?.data,
    isLoading: query.isLoading,
    isError: query.isError,
    error: query.error,
    refetch: query.refetch,
  }
}

export const useUserCreate = () => {
  const { invalidateUsersList } = useUsersCache()

  const mutation = useCreateUser({
    mutation: {
      onMutate: async () => {
        toast.loading('Creating user...', { id: 'create-user' })
        return {}
      },
      onSuccess: () => {
        invalidateUsersList()
        toast.success('User created successfully!', { id: 'create-user' })
      },
      onError: (error) => {
        toast.error(`Failed to create user: ${error.message}`, {
          id: 'create-user',
        })
      },
    },
  })

  return {
    createUser: (data: CreateUserRequest) =>
      mutation.mutateAsync({ data }),
    isCreating: mutation.isPending,
    isSuccess: mutation.isSuccess,
    isError: mutation.isError,
    error: mutation.error,
    reset: mutation.reset,
  }
}

export const useUserUpdate = () => {
  const { queryClient, invalidateUsersList } = useUsersCache()

  const mutation = useUpdateUser({
    mutation: {
      onMutate: async ({ id }) => {
        toast.loading('Updating user...', { id: `update-user-${id}` })

        await queryClient.cancelQueries({
          queryKey: getGetUserQueryKey(id),
        })

        const previousUser = queryClient.getQueryData(
          getGetUserQueryKey(id)
        )

        return { previousUser }
      },
      onSuccess: (_data, { id }) => {
        queryClient.invalidateQueries({
          queryKey: getGetUserQueryKey(id),
        })

        invalidateUsersList()

        toast.success('User updated successfully!', {
          id: `update-user-${id}`,
        })
      },
      onError: (error, { id }, context) => {
        if (context?.previousUser) {
          queryClient.setQueryData(
            getGetUserQueryKey(id),
            context.previousUser
          )
        }

        toast.error(`Failed to update user: ${error.message}`, {
          id: `update-user-${id}`,
        })
      },
    },
  })

  return {
    updateUser: (id: string, data: UpdateUserRequest) =>
      mutation.mutateAsync({ id, data }),
    isUpdating: mutation.isPending,
    isSuccess: mutation.isSuccess,
    isError: mutation.isError,
    error: mutation.error,
    reset: mutation.reset,
  }
}

export const useUserDelete = () => {
  const { queryClient, invalidateUsersList } = useUsersCache()

  const mutation = useDeleteUser({
    mutation: {
      onMutate: async ({ id }) => {
        toast.loading('Deleting user...', { id: `delete-user-${id}` })

        await queryClient.cancelQueries({
          queryKey: getListUsersQueryKey(),
          exact: false,
        })

        const previousUsers = queryClient.getQueryData(
          getListUsersQueryKey()
        )

        queryClient.setQueryData(
          getListUsersQueryKey(),
          (old: any) => {
            if (!old?.data) return old
            return {
              ...old,
              data: old.data.filter(
                (user: { id: string }) => user.id !== id
              ),
            }
          }
        )

        return { previousUsers }
      },
      onSuccess: (_, { id }) => {
        invalidateUsersList()

        queryClient.removeQueries({
          queryKey: getGetUserQueryKey(id),
        })

        toast.success('User deleted successfully!', {
          id: `delete-user-${id}`,
        })
      },
      onError: (error, { id }, context) => {
        if (context?.previousUsers) {
          queryClient.setQueryData(
            getListUsersQueryKey(),
            context.previousUsers
          )
        }

        toast.error(`Failed to delete user: ${error.message}`, {
          id: `delete-user-${id}`,
        })
      },
    },
  })

  return {
    deleteUser: (id: string) => mutation.mutateAsync({ id }),
    isDeleting: mutation.isPending,
    isSuccess: mutation.isSuccess,
    isError: mutation.isError,
    error: mutation.error,
    reset: mutation.reset,
  }
}
