export class Execution {
  constructor(
    public id: string,
    public algo: string,
    public data: string,
    public local_args: string[],
    public global_timeout: number,
    public local_timeout: number,
    public preprocess_args: string[],
    public mode: string,
    public global_args: string[],
    public files: string[],
    public model_token: string,
    public model_name: string,
    public state: string,
    public result: string,
  ) { }
}

export class ExecutionReq {
  constructor(
    public data: string,
    public local_args: string[],
    public type: string,
    public global_timeout: number,
    public local_timeout: number,
    public preprocess_args: string[],
    public mode: string,
    public global_args: string[],
    public files: string[],
  ) { }
}

export class StartExecutionReq {
  constructor(
    public algo: string,
    public executions: ExecutionReq[],
  ) { }
}
