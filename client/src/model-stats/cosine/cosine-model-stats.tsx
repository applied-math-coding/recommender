import { Space, Spin } from 'antd';
import { useEffect } from 'react';
import { changePageTitle, changeSelectedMenuKey } from '../../app.slice';
import { FetchState } from '../../domain/fetch-state';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { changeStatsLoading, fetchCosineModelStats } from './cosine-model-stats.slice';
import "./cosine-model-stats.scss";
import { Histogram } from '@ant-design/charts';
import { listToChartData } from '../../common/utils';
import { Typography } from 'antd';
import { MenuKey } from '../../domain/menu-key.enum';

const { Title } = Typography;

export default function CosineModelStats() {
  const PAGE_TITLE = 'Model Statistics';
  const statsLoading = useAppSelector(s => s.cosineModelStats.statsLoading);
  const cosineStatistic = useAppSelector(s => s.cosineModelStats.cosineStatistic);
  const pageTitle = useAppSelector(s => s.app.pageTitle);
  const dispatch = useAppDispatch();

  useEffect(() => {
    if (pageTitle !== PAGE_TITLE) {
      // this triggers a full re-render
      dispatch(changePageTitle(PAGE_TITLE));
      dispatch(changeSelectedMenuKey(MenuKey.ModelStatistics));
      dispatch(changeStatsLoading(FetchState.idle));
    } else if (statsLoading === FetchState.idle) {
      dispatch(fetchCosineModelStats());
    }
  });

  const noDataContainer = (
    <Title level={4} type="secondary">No Model found.</Title>
  )

  const loader = (
    <div className="loader-positioner">
      <Space size="middle">
        <Spin size="large" />
      </Space>
    </div>
  )

  const statsContainer = (
    <>
      <Title level={4} style={{ textAlign: 'center' }}>Distributions of Cosines</Title>
      {
        statsLoading === FetchState.succeeded &&
        <div className="chart-container">
          <Histogram binWidth={0.01} binField={'value'} data={listToChartData(cosineStatistic.cosines)} />
        </div>
      }
    </>
  )

  return (
    <>
      {
        statsLoading === FetchState.loading &&
        loader
      }
      {
        statsLoading === FetchState.succeeded && cosineStatistic?.cosines.length === 0 &&
        noDataContainer
      }
      {
        statsLoading === FetchState.succeeded && cosineStatistic?.cosines.length > 0 &&
        statsContainer
      }
    </>
  )
}
