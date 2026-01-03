import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/cek/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/cek/"!</div>
}
