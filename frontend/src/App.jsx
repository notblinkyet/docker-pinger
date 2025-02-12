import { useEffect, useState } from "react";
import { Table, Button, Input, message } from "antd";
import axios from "axios";

const backendUrl = import.meta.env.VITE_BACKEND_URL;
const refreshInterval = parseInt(import.meta.env.VITE_REFRESH_INTERVAL, 10)

const App = () => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [newIp, setNewIp] = useState("");

  const fetchData = async () => {
    setLoading(true);
    try {
      const response = await axios.get(`${backendUrl}/pings/last`);
      const processedData = response.data.data.map(item => ({
        ...item,
        last_success: item.was_success_before ? item.last_success : "",
      }));
      setData(processedData);
    } catch (error) {
      console.error("Ошибка загрузки данных:", error);
      message.error("Ошибка загрузки данных");
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, refreshInterval);
    return () => clearInterval(interval);
  }, []);

  const addIp = async () => {
    if (!newIp) return message.warning("Введите IP-адрес");
    try {
      await axios.post(`${backendUrl}/containers`, { ip: newIp });
      message.success("IP добавлен!");
      setNewIp("");
      fetchData();
    } catch (error) {
      console.error("Ошибка добавления IP:", error);
      message.error("Ошибка добавления IP");
    }
  };

  const deleteIp = async (ip) => {
    try {
      await axios.delete(`${backendUrl}/containers/${ip}`);
      message.success("IP удален!");
      fetchData();
    } catch (error) {
      console.error("Ошибка удаления IP:", error);
      message.error("Ошибка удаления IP");
    }
  };

  const columns = [
    { title: "IP-адрес", dataIndex: "ip", key: "ip" },
    { title: "Время пинга", dataIndex: "ping_at", key: "ping_at" },
    { title: "Последний успешный пинг", dataIndex: "last_success", key: "last_success" },
    { title: "Время отклика на пинг (мс)", dataIndex: "latency", key: "latency", render: (latency) => latency / 1000 },
    {
      title: "Действия",
      key: "actions",
      render: (_, record) => <Button danger onClick={() => deleteIp(record.ip)}>Удалить</Button>,
    },
  ];

  return (
    <div style={{ padding: "20px" }}>
      <h1>Мониторинг IP-адресов</h1>
      <div style={{ marginBottom: "20px", display: "flex", gap: "10px" }}>
        <Input placeholder="Введите IP" value={newIp} onChange={(e) => setNewIp(e.target.value)} />
        <Button type="primary" onClick={addIp}>Добавить IP</Button>
      </div>
      <Table dataSource={data} columns={columns} rowKey="ip" loading={loading} />
    </div>
  );
};

export default App;