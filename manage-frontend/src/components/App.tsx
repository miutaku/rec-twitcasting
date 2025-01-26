import React, { useEffect, useState } from 'react';
import { fetchCastingUsers } from '../services/api';
import '../styles/App.css'; // 修正されたパス

const App: React.FC = () => {
    const [users, setUsers] = useState([]);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            const result = await fetchCastingUsers();
            if (result.error) {
                setError(result.error);
            } else {
                setUsers(result);
            }
        };
        fetchData();
    }, []);

    return (
        <div className="App">
            <h1>キャスティングユーザー一覧</h1>
            {error && <p className="error">{error}</p>}
            <ul>
                {users.map((user: any) => (
                    <li key={user.target_username}>{user.target_username}</li>
                ))}
            </ul>
        </div>
    );
};

export default App;