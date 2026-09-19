import { createSlice } from "@reduxjs/toolkit";

interface Report {
  id: string;
  title: string;
  createdAt: string;
}

interface ReportsState {
  list: Report[];
  loading: boolean;
  error: string | null;
}

const initialState: ReportsState = {
  list: [],
  loading: false,
  error: null,
};

const reportsSlice = createSlice({
  name: "reports",
  initialState,
  reducers: {
    setReports: (state, action) => {
      state.list = action.payload;
    },
    setLoading: (state, action) => {
      state.loading = action.payload;
    },
    setError: (state, action) => {
      state.error = action.payload;
    },
  },
});

export const { setReports, setLoading, setError } = reportsSlice.actions;
export default reportsSlice.reducer;
