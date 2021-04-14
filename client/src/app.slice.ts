import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { MenuKey } from './domain/menu-key.enum';
import { ModelType } from './domain/model-type.enum';

export interface AppState {
  pageTitle: string;
  selectedMenuKey: MenuKey;
  selectedModelType: ModelType;
}

const initialState: AppState = {
  pageTitle: '',
  selectedMenuKey: MenuKey.DataImport,
  selectedModelType: ModelType.APRIORI
};

export const appSlice = createSlice({
  name: 'app',
  initialState,
  reducers: {
    changePageTitle: (state, action: PayloadAction<string>) => {
      state.pageTitle = action.payload;
    },
    changeSelectedMenuKey: (state, action: PayloadAction<MenuKey>) => {
      state.selectedMenuKey = action.payload;
    },
    changeSelectedModelType: (state, action: PayloadAction<ModelType>) => {
      state.selectedModelType = action.payload;
    }
  }
});

export const { changePageTitle, changeSelectedMenuKey, changeSelectedModelType } = appSlice.actions
export default appSlice.reducer;
