defmodule Example do
  use GenServer, restart: :transient

  @impl true
  def init(_opts) do
    IO.inspect("STARTING GENSERVER")
    Process.flag(:trap_exit, true)
    {:ok, %{}}
  end

  def start_link(opts) do
    GenServer.start_link(__MODULE__, opts, name: __MODULE__)
  end

  def terminate(reason, state) do
    IO.inspect("PROCESS TERMINATED: #{reason}")

    {:noreply, state}
  end

  @impl true
  def handle_info({:DOWN, _ref, :process, _pid, reason}, state) do
    IO.inspect("CHILD PROCESS IS DOWN: #{reason}")
    {:noreply, state}
  end

  @impl true
  def handle_call(:expire, _from, state) do
    {pid, _ref} =
      spawn_monitor(fn ->
        Process.flag(:trap_exit, true)

        # Process.send_after(self(), :msg, 1000)
        Process.send_after(self(), :ex, 1000)

        receive do
          :msg ->
            IO.inspect("RECEIVED A MESSAGE")

          :ex ->
            raise "yada yada"

          {:DOWN, _ref, :process, _pid, reason} ->
            IO.inspect("RECEIVED AN ERROR #{reason}")
            exit(reason)
        end
      end)

    {:reply, pid, state}
  end

  def expire(pid) do
    GenServer.call(pid, :expire)
  end
end

Supervisor.start_link([Example], strategy: :one_for_all)
