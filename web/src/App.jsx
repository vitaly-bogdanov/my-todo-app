import { useUsers } from './state/useUsers'
import { useTasks } from './state/useTasks'
import { useStatistics } from './state/useStatistics'
import { UsersView } from './view/UsersView'
import { TasksView } from './view/TasksView'
import { StatisticsView } from './view/StatisticsView'
import './view/styles.css'

export function App() {
  const users = useUsers()
  const tasks = useTasks()
  const statistics = useStatistics()

  return (
    <main>
      <h1>My Todo App</h1>
      <UsersView
        users={users.users}
        total={users.total}
        error={users.error}
        loading={users.loading}
        onLoad={users.load}
        onCreate={users.create}
        onGet={users.getOne}
        onPatch={users.patch}
        onRemove={users.remove}
      />
      <TasksView
        tasks={tasks.tasks}
        total={tasks.total}
        error={tasks.error}
        loading={tasks.loading}
        onLoad={tasks.load}
        onCreate={tasks.create}
        onGet={tasks.getOne}
        onPatch={tasks.patch}
        onRemove={tasks.remove}
      />
      <StatisticsView
        statistics={statistics.statistics}
        error={statistics.error}
        loading={statistics.loading}
        onLoad={statistics.load}
      />
    </main>
  )
}
