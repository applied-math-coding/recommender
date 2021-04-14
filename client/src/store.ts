import { configureStore } from '@reduxjs/toolkit'
import dataImportReducer from './data-import/data-import.slice';
import appReducer from './app.slice';
import modelStatsReducer from './model-stats/apriori/model-stats.slice';
import modelUseReducer from './model-use/apriori/model-use.slice';
import cosineModelStatsReducer from './model-stats/cosine/cosine-model-stats.slice';
import cosineModelUseReducer from './model-use/cosine/cosine-model-use.slice';

export const store = configureStore({
  reducer: {
    cosineModelStats: cosineModelStatsReducer,
    modelStats: modelStatsReducer,
    dataImport: dataImportReducer,
    modelUse: modelUseReducer,
    cosineModelUse: cosineModelUseReducer,
    app: appReducer
  }
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch