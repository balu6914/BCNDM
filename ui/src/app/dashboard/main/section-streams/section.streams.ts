export enum StreamsType {
  Bought = "Bought",
  Sold = "Sold"
}

export class StreamSection {
  name: StreamsType = StreamsType.Bought;
}
