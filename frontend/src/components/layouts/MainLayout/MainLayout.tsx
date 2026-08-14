import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';

import ActivityContainer from '../../Activity/ActivityContainer';
import ApplicationList from '../../Applications/ApplicationList';
import LatestInstanceStats from '../../Instances/LatestInstanceStats';

function MainLayout() {
  return (
    <Grid container spacing={2} justifyContent="center" alignItems="flex-start">
      <Grid
        size={{
          xs: 12,
          sm: 8,
        }}
      >
        <ApplicationList />
      </Grid>
      <Grid
        size={{
          xs: 12,
          sm: 4,
        }}
      >
        <ActivityContainer />
      </Grid>
      <Grid size={12}>
        <Paper sx={{ p: 2 }}>
          <LatestInstanceStats />
        </Paper>
      </Grid>
    </Grid>
  );
}

export default MainLayout;
