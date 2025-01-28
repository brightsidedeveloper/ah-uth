import { createLazyFileRoute } from '@tanstack/react-router'
import { useEffect } from 'react'
import { User, Users } from '../api/api'
import { get, post } from '../api/request'
import { useQuery } from '@tanstack/react-query'

export const Route = createLazyFileRoute('/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { data, error } = useQuery({
    queryKey: ['/api/users'],
    queryFn: () => get('/api/auth/gothic', undefined, (b) => User.decode(b)),
  })

  useEffect(() => {
    get('/api/users', undefined, (b) => Users.decode(b))
  }, [])

  return (
    <div>
      <button
        onClick={async () => {
          post(
            '/api/user',
            User.create({
              name: 'test2',
            }),
            (u) => User.encode(u).finish(),
            (b) => User.decode(b)
          )
        }}
      >
        Click
      </button>
      <br />
      {data?.name}
      <br />
      {error?.message}
      <br />
      {data?.name ? (
        <button
          onClick={() => {
            window.location.href = 'http://localhost:8081/api/auth/google/logout'
          }}
        >
          Logout
        </button>
      ) : (
        <button
          onClick={() => {
            window.location.href = 'http://localhost:8081/api/auth/google'
          }}
        >
          Login
        </button>
      )}
    </div>
  )
}
