import { useState, useEffect } from "react";

import "./App.css";
import axios from "./config/axios";

interface SentData {
	from: string;
	to: string;
	amount: number;
}

interface ModelData {
	from: string;
	to: string;
	rate: number;
	amount: number;
	result: number;
}

export default function App() {
	const [currencyList, setCurrencyList] = useState<string[]>([]);
	const [result, setResult] = useState<ModelData>();
	const [loading, setLoading] = useState<boolean>(false);

	const [sentData, setSentData] = useState<SentData>({
		from: "",
		to: "",
		amount: 0
	});

	const limit = 1000000000000000;

	const SetData = (key: string, value: string | number) => {
		setSentData(prev => ({ ...prev, [key]: value }))
	}

	const SendData = async (e: React.FormEvent) => {
		e.preventDefault();
		if (loading)
			return;

		if (!sentData.from || !sentData.to) {
			alert("All data must be filled.");
			return;
		}

		if (sentData.amount <= 0 || sentData.amount > limit) {
			alert(`Amount must be more than 0 and lesser than ${limit}.`);
			return;
		}

		setLoading(true);
		try {
			const res = await axios.post("/convert", sentData);
			console.log(res.data);
			setResult(res.data || null);
		} catch (err) {
			alert(err);
			return;
		} finally {
			setLoading(false);
		}
	}

	useEffect(() => {
		const GetCurrencies = async () => {
			try {
				const res = await axios.get("/get-currencies");
				setCurrencyList(res.data.list || []);
			} catch (err) {
				setCurrencyList([]);
				alert(err);
				return;
			}
		}

		GetCurrencies();
	})

	return (
		<div className="flex justify-center items-center min-h-screen mx-4">
			<div className="bg-white w-full max-w-lg rounded-xl shadow-2xl p-5 md:p-8 text-center transition-all duration-300">
                <h1 className="text-2xl md:text-3xl font-bold mb-5 text-gray-800">Currency Converter</h1>
                <hr className="border-gray-200" />
				<div className="pt-5">
					<form className="space-y-4" onSubmit={(e) => SendData(e)}>
						<div className="flex">
							<label htmlFor="amount" className="block text-left p-1">Amount:</label>
							<input min={0} max={limit} id="amount" type="number" placeholder="Amount" className="border-b p-1 w-full"
								onChange={(e) => SetData("amount", +e.target.value)} />
						</div>
						<div className="flex flex-col md:flex-row gap-4">
							<div className="w-full">
                                <label className="block text-left text-sm font-semibold text-gray-500 mb-1 ml-1">From</label>
                                <input type="text" value={sentData.from} placeholder="USD" 
                                    className="border border-gray-300 p-3 rounded-lg w-full text-center uppercase focus:outline-none focus:ring-2 focus:ring-blue-400"
                                    list="currencies" onChange={(e) => SetData("from", e.target.value)} />
                            </div>
                            <div className="hidden md:flex mt-4 items-center text-gray-400">➜</div>
                            <div className="w-full">
                                <label className="block text-left text-sm font-semibold text-gray-500 mb-1 ml-1">To</label>
                                <input type="text" value={sentData.to} placeholder="IDR" 
                                    className="border border-gray-300 p-3 rounded-lg w-full text-center uppercase focus:outline-none focus:ring-2 focus:ring-blue-400"
                                    list="currencies" onChange={(e) => SetData("to", e.target.value)} />
                            </div>
						</div>
						<datalist id="currencies">
							{
								currencyList.map((c, id) => (
									<option key={id} value={c} />
								))
							}
						</datalist>

						<button className="mt-3 p-2 rounded-xl text-xl bg-blue-400 w-full text-white font-bold hover:cursor-pointer">
							Convert
						</button>
					</form>

					{result && (
						<div className="mt-8 bg-gradient-to-r from-gray-700 to-gray-600 w-full rounded-xl shadow-inner p-6 text-center text-white animate-fade-in-up">
                            <p className="text-sm md:text-base tracking-wide break-all opacity-90 mb-1">
                                {result.amount.toLocaleString()} {result.from} =
                            </p>
                            <p className="text-3xl md:text-4xl tracking-wide break-all font-bold tracking-wide">
                                {result.result.toLocaleString(undefined, { maximumFractionDigits: 2 })} {result.to}
                            </p>
                            <p className="text-2xl text-wrap text-gray-300 mt-4">
                                Exchange Rate: 1 {result.from} ≈ {result.rate.toLocaleString()} {result.to}
                            </p>
                        </div>
					)}
				</div>
			</div>
		</div>
	);
}