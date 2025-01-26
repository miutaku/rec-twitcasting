import React, { useEffect, useState } from 'react';
import { fetchCastingUsers } from '../services/api';
import './App.css';

const App: React.FC = () => {
    const [users, setUsers] = useState<any[]>([]);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const getUsers = async () => {
            try {
                const data = await fetchCastingUsers();
                setUsers(data);
            } catch (err) {
                setError('データの取得に失敗しました');
            }
        };

        getUsers();
    }, []);

    return (
        <div className="App">
            <h1>キャスティングユーザー一覧</h1>
            {error && <p className="error">{error}</p>}
            <ul>
                {users.map((user) => (
                    <li key={user.target_username}>{user.target_username}</li>
                ))}
            </ul>
        </div>
    );
};

export default App;