/* SystemJS module definition */
declare var module: NodeModule;
interface NodeModule {
  id: string;
}

/* Map Wildcard Module */
declare module "map-style.json" {
    const value: any;
    export default value;
}
