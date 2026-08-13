import Box from '@mui/material/Box';
import { useTheme } from '@mui/material/styles';
import React from 'react';
import { useTranslation } from 'react-i18next';

import API from '../../api/API';
import { InstanceStats } from '../../api/apiDataTypes';
import Empty from '../common/EmptyContent';
import Loader from '../common/Loader';
import SimpleTable from '../common/SimpleTable';

const columns = {
  channel_name: 'Channel',
  arch: 'Architecture',
  version: 'Version',
  instances: 'Instances',
};

export default function LatestInstanceStats() {
  const [instanceStats, setInstanceStats] = React.useState<InstanceStats[] | null>(null);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState(false);
  const { t } = useTranslation();
  const theme = useTheme();

  React.useEffect(() => {
    API.getInstanceStatsLatest()
      .then(stats => {
        setInstanceStats(stats);
        setLoading(false);
      })
      .catch(err => {
        console.error('Error getting latest instance stats from the API:', err);
        setError(true);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <Loader />;
  }

  if (error) {
    return <Empty>{t('instances|latest_instance_stats_error')}</Empty>;
  }

  // Empty array is a valid response (no data yet), not an error.
  if (!instanceStats || instanceStats.length === 0) {
    return <Empty>{t('instances|latest_instance_stats_empty')}</Empty>;
  }

  return (
    <Box>
      <Box fontSize={18} fontWeight={700} color={theme.palette.greyShadeColor} mb={2}>
        {t('instances|latest_instance_stats')}
      </Box>
      <SimpleTable
        emptyMessage={t('instances|latest_instance_stats_empty')}
        columns={columns}
        instances={instanceStats}
      />
    </Box>
  );
}
