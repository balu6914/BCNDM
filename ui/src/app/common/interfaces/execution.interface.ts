export class Execution {
  constructor(
    public id: string,
    public state: string,
    public algo: string,
    public data: string,
    public mode: string,
  ) { }
}

export class ExecutionReq {
  constructor(
    public data: string,
    public mode: string,
  ) { }
}

export class StartExecutionReq {
  constructor(
    public algo: string,
    public executions: ExecutionReq[],
  ) { }
}
