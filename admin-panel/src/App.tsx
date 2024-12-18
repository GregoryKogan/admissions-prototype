import {
  Admin,
  Resource,
  ListGuesser,
  EditGuesser,
  ShowGuesser,
} from "react-admin";
import { Layout } from "./Layout";

import { fetchUtils } from "react-admin";
import postgrestRestProvider, {
  IDataProviderConfig,
  defaultPrimaryKeys,
  defaultSchema,
} from "@raphiniert/ra-data-postgrest";

export const config: IDataProviderConfig = {
  apiUrl: "http://localhost:3333",
  httpClient: fetchUtils.fetchJson,
  defaultListOp: "eq",
  primaryKeys: defaultPrimaryKeys,
  schema: defaultSchema,
};


export const App = () => (
  <Admin layout={Layout} dataProvider={postgrestRestProvider(config)}>
    <Resource name="users" list={ListGuesser} edit={EditGuesser} show={ShowGuesser} />
    <Resource name="registration_data" list={ListGuesser} edit={EditGuesser} show={ShowGuesser} />
    <Resource name="exams" list={ListGuesser} edit={EditGuesser} show={ShowGuesser} />
    <Resource name="exam_registrations" list={ListGuesser} edit={EditGuesser} show={ShowGuesser} />
    <Resource name="exam_results" list={ListGuesser} edit={EditGuesser} show={ShowGuesser} />
    <Resource name="exam_types" list={ListGuesser} edit={EditGuesser} show={ShowGuesser} />
    <Resource name="roles" list={ListGuesser} edit={EditGuesser} show={ShowGuesser} />
    <Resource name="passwords" list={ListGuesser} edit={EditGuesser} show={ShowGuesser} />
  </Admin>
);
