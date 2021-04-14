import { useEffect } from 'react';
import { changePageTitle, changeSelectedMenuKey } from '../../app.slice';
import { FetchState } from '../../domain/fetch-state';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { fetchCosineRecommendations } from './cosine-model-use.slice';
import { Space, Spin, Typography } from 'antd';
import Search from 'antd/lib/input/Search';
import { MenuKey } from '../../domain/menu-key.enum';
import './cosine-model-use.scss';

const { Title } = Typography;

export default function CosineModelUse() {
  const PAGE_TITLE = 'Model Usage';
  const dispatch = useAppDispatch();
  const loadingRecommendations = useAppSelector(state => state.cosineModelUse.loadingRecommendations);
  const pageTitle = useAppSelector(s => s.app.pageTitle);
  const recommendations = useAppSelector(state => state.cosineModelUse.recommendations);

  useEffect(() => {
    if (pageTitle !== PAGE_TITLE) {
      // this triggers a full re-render
      dispatch(changePageTitle('Model Usage'));
      dispatch(changeSelectedMenuKey(MenuKey.ModelUsage));
    }
  })

  function handleTest(v: string) {
    dispatch(fetchCosineRecommendations(v.split(',')));
  }

  return (
    <>
      <Title level={4} style={{ marginBottom: '10px' }}>Test Recommendations</Title>
      <Search
        placeholder="1234,5434,4563"
        allowClear
        enterButton="Send"
        size="large"
        onSearch={handleTest}
        style={{ width: '50%' }}
      />
      {
        loadingRecommendations === FetchState.loading &&
        <div className="loader-positioner">
          <Space size="middle">
            <Spin size="large" />
          </Space>
        </div>
      }
      {
        loadingRecommendations === FetchState.succeeded &&
        <div style={{ marginTop: '20px', marginLeft: '20px' }}>
          <Title level={5} style={{ marginBottom: '5px' }}>Recommended Items:</Title>
          {
            recommendations.map(r =>
              <div key={r.item} className="recommendation">
                <div>Item: {r.item}</div>
                <div>Cosine: {Math.round(r.cosine * 100)}%</div>
              </div>
            )
          }
        </div>
      }
    </>
  )
}
