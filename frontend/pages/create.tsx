import Layout from "@/components/Layout"
import PollCreation from "@/components/PollCreation"

export default function CreatePoll() {
  return (
    <Layout>
      <h1 className="text-3xl font-bold mb-6">Create a New Poll</h1>
      <PollCreation />
    </Layout>
  )
}

