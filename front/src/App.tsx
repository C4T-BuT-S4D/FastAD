import {useEffect, useState} from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import {Team} from "./proto/data/teams/teams.ts";
import {teamsService} from "./services/data";

function App() {
    const [count, setCount] = useState(0);

    const [teams, setTeams] = useState<Team[]>();

    useEffect(() => {
        const wrapper = async () => {
            const resp = await teamsService.list({lastUpdate: 0});
            console.log("fetched teams", resp);
            setTeams(resp.teams);
        }

        void wrapper();
    });

    return (
        <>
            <div>
                <a href="https://vitejs.dev" target="_blank">
                    <img src={viteLogo} className="logo" alt="Vite logo"/>
                </a>
                <a href="https://react.dev" target="_blank">
                    <img src={reactLogo} className="logo react" alt="React logo"/>
                </a>
            </div>
            <h1>Vite + React</h1>
            <div className="card">
                <button onClick={() => setCount((count) => count + 1)}>
                    count is {count}
                </button>
                <p>
                    Edit <code>src/App.tsx</code> and save to test HMR
                </p>
            </div>
            <h4>Teams:</h4>
            {teams?.map((team) => (
                <p>{team.name}</p>
            ))}
            <p className="read-the-docs">
                Click on the Vite and React logos to learn more
            </p>
        </>
    )
}

export default App
