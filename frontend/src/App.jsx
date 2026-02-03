import React, { useState, useEffect } from 'react';
import './App.css';

const MOCK_JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJlbWFpbCI6ImZvdW5kZXJAbGFnLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.NO_SIG";

function App() {
  const [agency, setAgency] = useState(null);
  const [loading, setLoading] = useState(true);
  const [showCreateForm, setShowCreateForm] = useState(false);

  const [dailySnapshot, setDailySnapshot] = useState(null);
  const [snapshotLoading, setSnapshotLoading] = useState(false);
  const [todayCashInput, setTodayCashInput] = useState('');
  const [snapshotError, setSnapshotError] = useState('');

  const [dailySummary, setDailySummary] = useState(null);
  const [burnRunway, setBurnRunway] = useState(null);
  const [retainerSummary, setRetainerSummary] = useState(null);
  const [clients, setClients] = useState([]);
  const [utilization, setUtilization] = useState(null);
  const [realityScore, setRealityScore] = useState(null);
  const [costBreakdown, setCostBreakdown] = useState(null);
  const [showBreakdown, setShowBreakdown] = useState(false);

  const [timeHours, setTimeHours] = useState('');
  const [timeClientId, setTimeClientId] = useState('');
  const [newClientName, setNewClientName] = useState('');
  const [selectedClientId, setSelectedClientId] = useState('');
  const [retainerAmount, setRetainerAmount] = useState('');
  const [revenueAmount, setRevenueAmount] = useState('');
  const [revenueSource, setRevenueSource] = useState('');
  const [revenueSubmitting, setRevenueSubmitting] = useState(false);
  const [costAmount, setCostAmount] = useState('');
  const [costType, setCostType] = useState('fixed');
  const [costLabel, setCostLabel] = useState('');
  const [costCategory, setCostCategory] = useState('other');
  const [costSubmitting, setCostSubmitting] = useState(false);
  const [name, setName] = useState('');
  const [currency, setCurrency] = useState('USD');
  const [startingCash, setStartingCash] = useState(0);

  useEffect(() => { fetchAgency(); }, []);

  const fetchAgency = async () => {
    setLoading(true);
    try {
      const res = await fetch('/api/agency', { headers: { 'CF-Access-Jwt-Assertion': MOCK_JWT } });
      if (res.status === 200) {
        const data = await res.json();
        setAgency(data);
        setShowCreateForm(false);
        fetchAll();
      } else if (res.status === 404) {
        setAgency(null);
        setShowCreateForm(true);
      }
    } catch (err) { console.error(err); }
    finally { setLoading(false); }
  };

  const fetchAll = () => {
    fetchDailySnapshot();
    fetchDailySummary();
    fetchBurnRunway();
    fetchRetainerSummary();
    fetchClients();
    fetchUtilization();
    fetchRealityScore();
    fetchCostBreakdown();
  };

  const fetchDailySnapshot = async () => {
    setSnapshotLoading(true);
    try {
      const res = await fetch('/api/cash-snapshot/today', { headers: { 'CF-Access-Jwt-Assertion': MOCK_JWT } });
      if (res.status === 200) setDailySnapshot(await res.json());
      else setDailySnapshot(null);
    } catch (err) { console.error(err); }
    finally { setSnapshotLoading(false); }
  };

  const fetchDailySummary = async () => {
    try {
      const res = await fetch('/api/daily-summary/today', { headers: { 'CF-Access-Jwt-Assertion': MOCK_JWT } });
      if (res.status === 200) setDailySummary(await res.json());
    } catch (err) { console.error(err); }
  };

  const fetchBurnRunway = async () => {
    try {
      const res = await fetch('/api/burn-runway', { headers: { 'CF-Access-Jwt-Assertion': MOCK_JWT } });
      if (res.status === 200) setBurnRunway(await res.json());
    } catch (err) { console.error(err); }
  };

  const fetchRetainerSummary = async () => {
    try {
      const res = await fetch('/api/retainer-summary', { headers: { 'CF-Access-Jwt-Assertion': MOCK_JWT } });
      if (res.status === 200) setRetainerSummary(await res.json());
    } catch (err) { console.error(err); }
  };

  const fetchClients = async () => {
    try {
      const res = await fetch('/api/clients', { headers: { 'CF-Access-Jwt-Assertion': MOCK_JWT } });
      if (res.status === 200) setClients(await res.json() || []);
    } catch (err) { console.error(err); }
  };

  const fetchUtilization = async () => {
    try {
      const res = await fetch('/api/utilization', { headers: { 'CF-Access-Jwt-Assertion': MOCK_JWT } });
      if (res.status === 200) setUtilization(await res.json());
    } catch (err) { console.error(err); }
  };

  const fetchRealityScore = async () => {
    try {
      const res = await fetch('/api/agency-reality-score', { headers: { 'CF-Access-Jwt-Assertion': MOCK_JWT } });
      if (res.status === 200) setRealityScore(await res.json());
    } catch (err) { console.error(err); }
  };

  const fetchCostBreakdown = async () => {
    try {
      const res = await fetch('/api/cost-breakdown', { headers: { 'CF-Access-Jwt-Assertion': MOCK_JWT } });
      if (res.status === 200) setCostBreakdown(await res.json());
    } catch (err) { console.error(err); }
  };

  const handleCreateAgency = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch('/api/agency', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'CF-Access-Jwt-Assertion': MOCK_JWT },
        body: JSON.stringify({ name, base_currency: currency, starting_cash: Number(startingCash) })
      });
      if (res.status === 201) fetchAgency();
    } catch (err) { console.error(err); }
  };

  const handleSaveSnapshot = async (e) => {
    e.preventDefault();
    setSnapshotError('');
    try {
      const res = await fetch('/api/cash-snapshot', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'CF-Access-Jwt-Assertion': MOCK_JWT },
        body: JSON.stringify({ cash_balance: Number(todayCashInput) })
      });
      if (res.status === 201) fetchAll();
      else if (res.status === 409) setSnapshotError("Already logged.");
    } catch (err) { setSnapshotError("Error"); }
  };

  const handleAddRevenue = async (e) => {
    e.preventDefault();
    setRevenueSubmitting(true);
    try {
      const res = await fetch('/api/revenue', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'CF-Access-Jwt-Assertion': MOCK_JWT },
        body: JSON.stringify({ amount: Number(revenueAmount), source: revenueSource })
      });
      if (res.status === 201) { setRevenueAmount(''); setRevenueSource(''); fetchAll(); }
    } catch (err) { console.error(err); }
    finally { setRevenueSubmitting(false); }
  };

  const handleAddCost = async (e) => {
    e.preventDefault();
    setCostSubmitting(true);
    try {
      const res = await fetch('/api/cost', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'CF-Access-Jwt-Assertion': MOCK_JWT },
        body: JSON.stringify({ amount: Number(costAmount), type: costType, label: costLabel, category: costCategory })
      });
      if (res.status === 201) { setCostAmount(''); setCostLabel(''); setCostCategory('other'); fetchAll(); }
    } catch (err) { console.error(err); }
    finally { setCostSubmitting(false); }
  };

  const handleAddClient = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch('/api/clients', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'CF-Access-Jwt-Assertion': MOCK_JWT },
        body: JSON.stringify({ name: newClientName })
      });
      if (res.status === 201) { setNewClientName(''); fetchClients(); }
    } catch (err) { console.error(err); }
  };

  const handleAddRetainer = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch('/api/retainers', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'CF-Access-Jwt-Assertion': MOCK_JWT },
        body: JSON.stringify({ client_id: selectedClientId, monthly_amount: Number(retainerAmount) })
      });
      if (res.status === 201) { setRetainerAmount(''); setSelectedClientId(''); fetchAll(); }
    } catch (err) { console.error(err); }
  };

  const handleAddTimeEntry = async (e) => {
    e.preventDefault();
    try {
      const body = { hours: Number(timeHours) };
      if (timeClientId) body.client_id = timeClientId;
      const res = await fetch('/api/time-entry', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'CF-Access-Jwt-Assertion': MOCK_JWT },
        body: JSON.stringify(body)
      });
      if (res.status === 201) { setTimeHours(''); setTimeClientId(''); fetchAll(); }
    } catch (err) { console.error(err); }
  };

  const getScoreColor = (status) => {
    if (status === 'Healthy') return '#22c55e';
    if (status === 'Watch') return '#eab308';
    if (status === 'At Risk') return '#f97316';
    return '#ef4444';
  };

  const getMarginColor = (m) => m >= 0 ? '#22c55e' : '#ef4444';
  const getRunwayColor = (r) => r === null ? '#666' : r >= 6 ? '#22c55e' : r >= 3 ? '#eab308' : '#ef4444';
  const getCoverageColor = (c) => c >= 1.5 ? '#22c55e' : c >= 1 ? '#eab308' : '#ef4444';
  const getConcentrationColor = (p) => p > 0.5 ? '#ef4444' : p >= 0.3 ? '#eab308' : '#22c55e';
  const getUtilizationColor = (u) => u >= 100 ? '#ef4444' : u >= 85 ? '#eab308' : u >= 60 ? '#22c55e' : '#666';

  if (loading) return <div style={{ padding: '2rem', textAlign: 'center' }}>Loading...</div>;

  if (showCreateForm) {
    return (
      <div style={{ maxWidth: 400, margin: '4rem auto', padding: '2rem' }}>
        <h1 style={{ marginBottom: '2rem' }}>Create Your Agency</h1>
        <form onSubmit={handleCreateAgency}>
          <div style={{ marginBottom: '1rem' }}>
            <input type="text" placeholder="Agency Name" value={name} onChange={(e) => setName(e.target.value)} required style={{ width: '100%', padding: '0.75rem', fontSize: '1rem' }} />
          </div>
          <div style={{ marginBottom: '1rem' }}>
            <input type="text" placeholder="Currency (USD)" value={currency} onChange={(e) => setCurrency(e.target.value)} required style={{ width: '100%', padding: '0.75rem', fontSize: '1rem' }} />
          </div>
          <div style={{ marginBottom: '1rem' }}>
            <input type="number" placeholder="Starting Cash" value={startingCash} onChange={(e) => setStartingCash(e.target.value)} required style={{ width: '100%', padding: '0.75rem', fontSize: '1rem' }} />
          </div>
          <button type="submit" style={{ width: '100%', padding: '1rem', fontSize: '1rem', background: '#111', color: 'white', border: 'none', cursor: 'pointer' }}>Launch Agency</button>
        </form>
      </div>
    );
  }

  if (agency) {
    return (
      <div style={{ maxWidth: 800, margin: '0 auto', padding: '1rem', fontFamily: 'system-ui, sans-serif' }}>

        {/* HEADER: Score + Cash + Retainers */}
        {realityScore && (
          <div style={{ background: '#111', color: 'white', padding: '1.5rem', marginBottom: '2rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1rem' }}>
              <div>
                <div style={{ opacity: 0.6, fontSize: '0.85rem', marginBottom: '0.25rem' }}>Agency Reality Score</div>
                <div style={{ fontSize: '4rem', fontWeight: '700', lineHeight: 1, color: getScoreColor(realityScore.status) }}>
                  {realityScore.score}
                </div>
              </div>
              <div style={{ textAlign: 'right' }}>
                <div style={{ fontSize: '1.5rem', fontWeight: '600', color: getScoreColor(realityScore.status) }}>
                  {realityScore.status}
                </div>
                <div style={{ opacity: 0.5, fontSize: '0.8rem', cursor: 'pointer' }} onClick={() => setShowBreakdown(!showBreakdown)}>
                  {showBreakdown ? '▲ Hide' : '▼ Show'} breakdown
                </div>
              </div>
            </div>

            {/* Always visible: Cash + Retainers */}
            <div style={{ display: 'flex', gap: '2rem', paddingTop: '1rem', borderTop: '1px solid #333' }}>
              <div>
                <div style={{ opacity: 0.6, fontSize: '0.75rem' }}>Cash on Hand</div>
                <div style={{ fontSize: '1.25rem', fontWeight: '600' }}>{Math.round(realityScore.cash_on_hand).toLocaleString()}</div>
              </div>
              <div>
                <div style={{ opacity: 0.6, fontSize: '0.75rem' }}>Committed Retainers</div>
                <div style={{ fontSize: '1.25rem', fontWeight: '600' }}>{Math.round(realityScore.committed_retainers).toLocaleString()}/mo</div>
              </div>
            </div>

            {showBreakdown && (
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(5, 1fr)', gap: '1rem', marginTop: '1rem', paddingTop: '1rem', borderTop: '1px solid #333', fontSize: '0.85rem' }}>
                <div><div style={{ opacity: 0.6 }}>Retainer</div><div style={{ fontWeight: '600' }}>{realityScore.breakdown.retainer_safety}/25</div></div>
                <div><div style={{ opacity: 0.6 }}>Runway</div><div style={{ fontWeight: '600' }}>{realityScore.breakdown.runway}/20</div></div>
                <div><div style={{ opacity: 0.6 }}>Concentration</div><div style={{ fontWeight: '600' }}>{realityScore.breakdown.client_concentration}/20</div></div>
                <div><div style={{ opacity: 0.6 }}>Profit</div><div style={{ fontWeight: '600' }}>{realityScore.breakdown.profitability}/20</div></div>
                <div><div style={{ opacity: 0.6 }}>Capacity</div><div style={{ fontWeight: '600' }}>{realityScore.breakdown.capacity_pressure}/15</div></div>
              </div>
            )}
          </div>
        )}

        {/* SURVIVAL: Will I die? */}
        <div style={{ marginBottom: '2rem' }}>
          <div style={{ fontSize: '0.75rem', fontWeight: '600', textTransform: 'uppercase', opacity: 0.5, marginBottom: '0.75rem' }}>Survival</div>

          {!dailySnapshot && !snapshotLoading && (
            <div style={{ background: '#fef3c7', padding: '1rem', marginBottom: '1rem' }}>
              <div style={{ fontWeight: '600', marginBottom: '0.5rem' }}>Log today's cash balance</div>
              <form onSubmit={handleSaveSnapshot} style={{ display: 'flex', gap: '0.5rem' }}>
                <input type="number" placeholder="Cash balance" value={todayCashInput} onChange={(e) => setTodayCashInput(e.target.value)} required style={{ flex: 1, padding: '0.5rem' }} />
                <button type="submit" style={{ padding: '0.5rem 1rem', background: '#111', color: 'white', border: 'none' }}>Save</button>
              </form>
              {snapshotError && <div style={{ color: '#ef4444', marginTop: '0.5rem' }}>{snapshotError}</div>}
            </div>
          )}

          {burnRunway && (
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
              <div style={{ padding: '1rem', background: '#f8f8f8' }}>
                <div style={{ opacity: 0.6, fontSize: '0.85rem' }}>Panic Runway</div>
                <div style={{ fontSize: '2rem', fontWeight: '700', color: getRunwayColor(burnRunway.runway_months) }}>
                  {burnRunway.runway_months !== null ? `${burnRunway.runway_months} mo` : '∞'}
                </div>
                <div style={{ opacity: 0.5, fontSize: '0.75rem' }}>Assumes zero future revenue</div>
              </div>
              <div style={{ padding: '1rem', background: '#f8f8f8' }}>
                <div style={{ opacity: 0.6, fontSize: '0.85rem' }}>Monthly Burn</div>
                <div style={{ fontSize: '2rem', fontWeight: '700' }}>{Math.round(burnRunway.monthly_burn).toLocaleString()}</div>
                <div style={{ opacity: 0.5, fontSize: '0.75rem' }}>Fixed costs (30 days)</div>
              </div>
            </div>
          )}
        </div>

        {/* STABILITY: Am I structurally okay? */}
        <div style={{ marginBottom: '2rem' }}>
          <div style={{ fontSize: '0.75rem', fontWeight: '600', textTransform: 'uppercase', opacity: 0.5, marginBottom: '0.75rem' }}>Stability</div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '1rem' }}>
            {burnRunway && (
              <div style={{ padding: '1rem', background: '#f8f8f8' }}>
                <div style={{ opacity: 0.6, fontSize: '0.85rem' }}>Operating Margin</div>
                <div style={{ fontSize: '1.5rem', fontWeight: '700', color: getMarginColor(burnRunway.operating_margin) }}>
                  {burnRunway.operating_margin >= 0 ? '+' : ''}{Math.round(burnRunway.operating_margin).toLocaleString()}
                </div>
                <div style={{ opacity: 0.5, fontSize: '0.75rem' }}>Retainers − Fixed Costs</div>
              </div>
            )}
            {retainerSummary && (
              <>
                <div style={{ padding: '1rem', background: '#f8f8f8' }}>
                  <div style={{ opacity: 0.6, fontSize: '0.85rem' }}>Retainer Coverage</div>
                  <div style={{ fontSize: '1.5rem', fontWeight: '700', color: getCoverageColor(retainerSummary.coverage_ratio) }}>
                    {retainerSummary.coverage_ratio}x
                  </div>
                  <div style={{ opacity: 0.5, fontSize: '0.75rem' }}>vs fixed costs</div>
                </div>
                <div style={{ padding: '1rem', background: '#f8f8f8' }}>
                  <div style={{ opacity: 0.6, fontSize: '0.85rem' }}>Top Client Risk</div>
                  <div style={{ fontSize: '1.5rem', fontWeight: '700', color: getConcentrationColor(retainerSummary.top_client_percentage) }}>
                    {Math.round(retainerSummary.top_client_percentage * 100)}%
                  </div>
                  <div style={{ opacity: 0.5, fontSize: '0.75rem' }}>concentration</div>
                </div>
              </>
            )}
          </div>
        </div>

        {/* DRAG: What's killing the agency? (Ticket 11) */}
        {realityScore && realityScore.primary_risk !== 'Healthy' && (
          <div style={{ marginBottom: '2rem' }}>
            <div style={{ fontSize: '0.75rem', fontWeight: '600', textTransform: 'uppercase', opacity: 0.5, marginBottom: '0.75rem' }}>What's Dragging You Down</div>
            <div style={{ padding: '1rem', background: '#f8f8f8', borderLeft: '4px solid #333' }}>
              <ul style={{ listStyle: 'none', padding: 0, margin: 0, fontSize: '0.9rem' }}>
                <li style={{ marginBottom: '0.5rem' }}>• Primary drag: <strong>{realityScore.primary_risk}</strong>
                  {costBreakdown && realityScore.primary_risk === 'High Fixed Costs' && (
                    <span> ({costBreakdown.primary_driver.category} costs @ {Math.round(costBreakdown.primary_driver.amount).toLocaleString()} / mo)</span>
                  )}
                </li>
                {burnRunway && burnRunway.operating_margin < 0 && (
                  <li>• Structural loss: <strong style={{ color: '#ef4444' }}>–{Math.abs(Math.round(burnRunway.operating_margin)).toLocaleString()}</strong> / month</li>
                )}
              </ul>
            </div>
          </div>
        )}

        {/* PRESSURE: Am I overloaded? */}
        {utilization && (
          <div style={{ marginBottom: '2rem' }}>
            <div style={{ fontSize: '0.75rem', fontWeight: '600', textTransform: 'uppercase', opacity: 0.5, marginBottom: '0.75rem' }}>Pressure</div>
            <div style={{ padding: '1rem', background: '#f8f8f8' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <div style={{ opacity: 0.6, fontSize: '0.85rem' }}>Capacity Utilization</div>
                  <div style={{ fontSize: '1.5rem', fontWeight: '700', color: getUtilizationColor(utilization.utilization_percent) }}>
                    {utilization.utilization_percent}%
                  </div>
                </div>
                <div style={{ opacity: 0.5, fontSize: '0.85rem' }}>
                  {utilization.used_hours}h / {utilization.capacity_hours}h
                </div>
              </div>
            </div>
          </div>
        )}

        {/* TODAY: Did I bleed or heal? */}
        {dailySummary && (
          <div style={{ marginBottom: '2rem' }}>
            <div style={{ fontSize: '0.75rem', fontWeight: '600', textTransform: 'uppercase', opacity: 0.5, marginBottom: '0.75rem' }}>Today</div>
            <div style={{ padding: '1.5rem', background: dailySummary.net >= 0 ? '#f0fdf4' : '#fef2f2' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <span style={{ marginRight: '1.5rem' }}>Revenue: <strong>{dailySummary.revenue.toLocaleString()}</strong></span>
                  <span>Costs: <strong>{dailySummary.costs.toLocaleString()}</strong></span>
                </div>
                <div style={{ fontSize: '1.5rem', fontWeight: '700', color: dailySummary.net >= 0 ? '#22c55e' : '#ef4444' }}>
                  {dailySummary.net >= 0 ? '+' : ''}{dailySummary.net.toLocaleString()}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* INPUT FORMS */}
        <div style={{ borderTop: '1px solid #eee', paddingTop: '2rem' }}>
          <div style={{ fontSize: '0.75rem', fontWeight: '600', textTransform: 'uppercase', opacity: 0.5, marginBottom: '1rem' }}>Quick Actions</div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem', marginBottom: '1rem' }}>
            <form onSubmit={handleAddRevenue} style={{ display: 'flex', gap: '0.5rem' }}>
              <input type="number" placeholder="Revenue" value={revenueAmount} onChange={(e) => setRevenueAmount(e.target.value)} required style={{ width: '80px', padding: '0.5rem' }} />
              <input type="text" placeholder="Source" value={revenueSource} onChange={(e) => setRevenueSource(e.target.value)} required style={{ flex: 1, padding: '0.5rem' }} />
              <button type="submit" disabled={revenueSubmitting} style={{ padding: '0.5rem 1rem', background: '#22c55e', color: 'white', border: 'none' }}>+Rev</button>
            </form>
            <form onSubmit={handleAddCost} style={{ display: 'flex', gap: '0.5rem' }}>
              <input type="number" placeholder="Cost" value={costAmount} onChange={(e) => setCostAmount(e.target.value)} required style={{ width: '80px', padding: '0.5rem' }} />
              <select value={costType} onChange={(e) => setCostType(e.target.value)} style={{ padding: '0.5rem', width: '85px' }}>
                <option value="fixed">Fixed</option>
                <option value="variable">Variable</option>
              </select>
              <select value={costCategory} onChange={(e) => setCostCategory(e.target.value)} style={{ padding: '0.5rem', width: '85px' }}>
                <option value="people">People</option>
                <option value="tools">Tools</option>
                <option value="other">Other</option>
              </select>
              <input type="text" placeholder="Label" value={costLabel} onChange={(e) => setCostLabel(e.target.value)} required style={{ flex: 1, padding: '0.5rem', minWidth: '50px' }} />
              <button type="submit" disabled={costSubmitting} style={{ padding: '0.5rem 1rem', background: '#ef4444', color: 'white', border: 'none' }}>+Cost</button>
            </form>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '1rem' }}>
            <form onSubmit={handleAddTimeEntry} style={{ display: 'flex', gap: '0.5rem' }}>
              <input type="number" placeholder="Hours" value={timeHours} onChange={(e) => setTimeHours(e.target.value)} required style={{ width: '60px', padding: '0.5rem' }} step="0.5" />
              <select value={timeClientId} onChange={(e) => setTimeClientId(e.target.value)} style={{ flex: 1, padding: '0.5rem' }}>
                <option value="">Internal</option>
                {clients.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
              </select>
              <button type="submit" style={{ padding: '0.5rem 1rem', background: '#111', color: 'white', border: 'none' }}>Log</button>
            </form>
            <form onSubmit={handleAddClient} style={{ display: 'flex', gap: '0.5rem' }}>
              <input type="text" placeholder="Client name" value={newClientName} onChange={(e) => setNewClientName(e.target.value)} required style={{ flex: 1, padding: '0.5rem' }} />
              <button type="submit" style={{ padding: '0.5rem 1rem', background: '#111', color: 'white', border: 'none' }}>+Client</button>
            </form>
            <form onSubmit={handleAddRetainer} style={{ display: 'flex', gap: '0.5rem' }}>
              <select value={selectedClientId} onChange={(e) => setSelectedClientId(e.target.value)} required style={{ flex: 1, padding: '0.5rem' }}>
                <option value="">Client...</option>
                {clients.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
              </select>
              <input type="number" placeholder="Monthly" value={retainerAmount} onChange={(e) => setRetainerAmount(e.target.value)} required style={{ width: '80px', padding: '0.5rem' }} />
              <button type="submit" style={{ padding: '0.5rem 1rem', background: '#111', color: 'white', border: 'none' }}>+Ret</button>
            </form>
          </div>
        </div>
      </div>
    );
  }

  return <div>Something went wrong.</div>;
}

export default App;
