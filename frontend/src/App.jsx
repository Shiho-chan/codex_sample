import './App.css'

function App() {
  return (
    <div className="portfolio">
      <header>
        <h1>Shiho-chan's Portfolio</h1>
        <p>ようこそ、私のポートフォリオサイトへ！</p>
      </header>
      <section>
        <h2>自己紹介</h2>
        <p>ここに自己紹介文を入れます。興味のある分野や経歴などを書いてください。</p>
      </section>
      <section>
        <h2>作品</h2>
        <ul>
          <li>Project A - 紹介文など</li>
          <li>Project B - 紹介文など</li>
          <li>Project C - 紹介文など</li>
        </ul>
      </section>
      <section>
        <h2>連絡先</h2>
        <p>
          ご連絡は <a href="mailto:shihochan@example.com">shihochan@example.com</a> までお願いします。
        </p>
      </section>
    </div>
  )
}

export default App
