import { render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import API from '../../api/API';
import { InstanceStats } from '../../api/apiDataTypes';
import LatestInstanceStats from '../../components/Instances/LatestInstanceStats';
import '../../i18n/config';

vi.mock('../../api/API', () => ({
  default: {
    getInstanceStatsLatest: vi.fn(),
  },
}));

const mockStats: InstanceStats[] = [
  {
    timestamp: '2026-08-10T14:32:00Z',
    channel_name: 'stable',
    arch: 'AMD64',
    version: '4500.1.0',
    instances: 128,
  },
  {
    timestamp: '2026-08-10T14:32:00Z',
    channel_name: 'beta',
    arch: 'ARM64',
    version: '4500.0.0',
    instances: 12,
  },
];

describe('LatestInstanceStats', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders rows given mock data', async () => {
    vi.mocked(API.getInstanceStatsLatest).mockResolvedValue(mockStats);

    render(<LatestInstanceStats />);

    await waitFor(() => {
      expect(screen.getByText('stable')).toBeTruthy();
    });

    expect(screen.getByText('Channel')).toBeTruthy();
    expect(screen.getByText('Architecture')).toBeTruthy();
    expect(screen.getByText('Version')).toBeTruthy();
    expect(screen.getByText('Instances')).toBeTruthy();
    expect(screen.getByText('AMD64')).toBeTruthy();
    expect(screen.getByText('4500.1.0')).toBeTruthy();
    expect(screen.getByText('128')).toBeTruthy();
    expect(screen.getByText('beta')).toBeTruthy();
    expect(screen.getByText('ARM64')).toBeTruthy();
    expect(screen.getByText('12')).toBeTruthy();
  });

  it('renders empty state given an empty array', async () => {
    vi.mocked(API.getInstanceStatsLatest).mockResolvedValue([]);

    render(<LatestInstanceStats />);

    await waitFor(() => {
      expect(screen.getByTestId('empty')).toBeTruthy();
    });

    expect(screen.getByText('No instance stats available yet.')).toBeTruthy();
    expect(screen.queryByText('Channel')).toBeNull();
  });
});
