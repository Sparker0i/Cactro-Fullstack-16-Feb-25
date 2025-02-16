import Layout from "@/components/Layout"
import PollList from "@/components/PollList"

export default function Home() {
  return (
    <Layout>
      <h1 className="text-3xl font-bold mb-6">Quick Polling App</h1>
      <PollList />
    </Layout>
  )
}

