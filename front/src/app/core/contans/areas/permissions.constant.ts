export const ACTION_TYPES = {
  CREATE: 'CREATE',
  READ: 'READ',
  UPDATE: 'UPDATE',
  DELETE: 'DELETE',
  TRANSFER: 'TRANSFER'
} as const;

// Creamos un tipo a partir del objeto
export type ActionType = typeof ACTION_TYPES[keyof typeof ACTION_TYPES];

export interface ModulePermission {
  name: string;
  user_id: number;
  components: Component[];
}

export interface Component {
  id: string;
  name: string;
  actions: ActionComponent[];
}

export interface ActionComponent {
  id: string;
  type: string;
  name: string;
  component_id: string;
}

export const MODULE_PERMISSIONS: ModulePermission[] = [
  {
    name: 'enfermería',
    user_id: 36,
    components: [
      {
        id: '5714cba8-ba66-4257-93a2-0c769d3d17a9',
        name: 'atenciones',
        actions: [
          {
            id: '969aed94-199f-46e4-96c8-22052b8ec08f',
            type: ACTION_TYPES.CREATE,
            name: 'Crear atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: '914a1c5a-e985-45d9-a026-0775c9269482',
            type: ACTION_TYPES.UPDATE,
            name: 'Editar atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: '2c385676-c34f-4d62-a97a-915cf7e82727',
            type: ACTION_TYPES.READ,
            name: 'Ver atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: '052256a5-d533-4d9a-b9bb-13852bdb2a2f',
            type: ACTION_TYPES.TRANSFER,
            name: 'Transferir atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: 'ec07c2bb-790d-4856-b9dd-a026b23f5f66',
            type: ACTION_TYPES.DELETE,
            name: 'Eliminar atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
        ]
      }
    ]
  },
  {
    name: 'enfermería',
    user_id: 37,
    components: [
      {
        id: '5714cba8-ba66-4257-93a2-0c769d3d17a9',
        name: 'atenciones',
        actions: [
          {
            id: '969aed94-199f-46e4-96c8-22052b8ec08f',
            type: ACTION_TYPES.CREATE,
            name: 'Crear atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: '2c385676-c34f-4d62-a97a-915cf7e82727',
            type: ACTION_TYPES.READ,
            name: 'Ver atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: '052256a5-d533-4d9a-b9bb-13852bdb2a2f',
            type: ACTION_TYPES.TRANSFER,
            name: 'Transferir atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: 'ec07c2bb-790d-4856-b9dd-a026b23f5f66',
            type: ACTION_TYPES.DELETE,
            name: 'Eliminar atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
        ]
      }
    ]
  },
  {
    name: 'odontología',
    user_id: 37,
    components: [
      {
        id: '5714cba8-ba66-4257-93a2-0c769d3d17a9',
        name: 'atenciones',
        actions: [
          {
            id: '969aed94-199f-46e4-96c8-22052b8ec08f',
            type: ACTION_TYPES.CREATE,
            name: 'Crear atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: '2c385676-c34f-4d62-a97a-915cf7e82727',
            type: ACTION_TYPES.READ,
            name: 'Ver atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: '052256a5-d533-4d9a-b9bb-13852bdb2a2f',
            type: ACTION_TYPES.TRANSFER,
            name: 'Transferir atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
          {
            id: 'ec07c2bb-790d-4856-b9dd-a026b23f5f66',
            type: ACTION_TYPES.DELETE,
            name: 'Eliminar atención',
            component_id: '5714cba8-ba66-4257-93a2-0c769d3d17a9'
          },
        ]
      }
    ]
  }
]
